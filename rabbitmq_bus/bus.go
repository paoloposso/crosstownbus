package rabbitmqbus

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/paoloposso/crosstownbus/core"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Uri string
}

type EventBus struct {
	channel *amqp.Channel
}

func CreateEventBus(config RabbitMQConfig) (core.EventBus, error) {
	conn, err := amqp.Dial(config.Uri)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return EventBus{
		channel: channel,
	}, nil
}

func (pub EventBus) Publish(message interface{}) error {
	eventName := reflect.TypeOf(message).Name()
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ch := pub.channel
	err = ch.ExchangeDeclare(eventName, "fanout", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.Publish(
		eventName, // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	return nil
}

func (bus EventBus) Subscribe(event reflect.Type, eventHandler core.EventHandler, retryOptions *core.RetryOptions) error {
	eventName := event.Name()
	ch := bus.channel
	queueName := fmt.Sprintf("%s.%s_queue", eventName, reflect.TypeOf(eventHandler).Name())

	err := ch.ExchangeDeclare(eventName, "fanout", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var retryOptionsTable amqp.Table
	retryExchange := ""
	retryKey := ""

	if retryOptions != nil && retryOptions.MaxRetryTimes > 0 {
		deadLetterExchange, deadLetterRoutingKey, err := bus.createDeadLetter(queueName, retryOptions.RetrySeconds)
		if err != nil {
			return err
		}
		retryOptionsTable = amqp.Table{
			"x-dead-letter-exchange":    deadLetterExchange,
			"x-dead-letter-routing-key": deadLetterRoutingKey,
		}
		retryExchange = deadLetterExchange
		retryKey = deadLetterRoutingKey
	}

	queue, err := ch.QueueDeclare(queueName, false, true, false, false, retryOptionsTable)
	if err != nil {
		return err
	}

	if retryOptions != nil {
		if err := ch.QueueBind(queue.Name, queue.Name, retryExchange, false, nil); err != nil {
			return err
		}
	}
	if err = ch.QueueBind(queue.Name, "", eventName, false, nil); err != nil {
		return err
	}
	msgs, err := ch.Consume(
		queueName,
		"",
		retryOptions == nil,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if retryOptions == nil || retryOptions.MaxRetryTimes == 0 {
		go func() {
			for msg := range msgs {
				go eventHandler.Handle(msg.Body)
			}
		}()
	} else {
		bus.handleWithRetry(msgs, eventHandler, retryOptions.MaxRetryTimes, retryExchange, retryKey)
	}
	return nil
}

func (bus EventBus) handleWithRetry(msgs <-chan amqp.Delivery, eventHandler core.EventHandler, maxRetry int32, retryExchange, retryKey string) {
	go func() {
		for msg := range msgs {
			if err := eventHandler.Handle(msg.Body); err != nil {
				log.Printf("Error handling message: %s", err)
				bus.rejectMessage(msg, maxRetry, retryExchange, retryKey)
			} else {
				msg.Ack(false)
			}
		}
	}()
}

func (bus EventBus) rejectMessage(msg amqp.Delivery, maxRetry int32, retryExchange, retryKey string) {
	var redeliveryCount int32 = 0

	if msg.Headers == nil {
		msg.Headers = make(amqp.Table)
	}
	if _, exists := msg.Headers["x-redelivered-count"]; !exists {
		msg.Headers["x-redelivered-count"] = int32(0)
	}

	redeliveryCount = msg.Headers["x-redelivered-count"].(int32) + 1

	if redeliveryCount > maxRetry {
		msg.Ack(false)
	} else {
		msg.Headers["x-redelivered-count"] = redeliveryCount
		msg.Ack(false)
		bus.channel.Publish(retryExchange, retryKey, false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.Body),
				Headers:     msg.Headers,
			})
	}
}

func (bus EventBus) createDeadLetter(queueName string, retrySeconds uint32) (string, string, error) {
	exchangeNameDl := fmt.Sprintf("%s_retryexh", queueName)
	queueNameDl := fmt.Sprintf("%s_retry", queueName)

	if err := bus.channel.ExchangeDeclare(
		exchangeNameDl,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return "", "", err
	}

	options := amqp.Table{
		"x-dead-letter-exchange":    exchangeNameDl,
		"x-dead-letter-routing-key": queueName,
		"x-message-ttl":             int32(retrySeconds * 1000),
	}

	if _, err := bus.channel.QueueDeclare(queueNameDl, true, false, false, false, options); err != nil {
		return "", "", err
	}

	if err := bus.channel.QueueBind(
		queueNameDl,
		queueNameDl,
		exchangeNameDl,
		false,
		nil,
	); err != nil {
		return "", "", err
	}

	return exchangeNameDl, queueNameDl, nil
}
