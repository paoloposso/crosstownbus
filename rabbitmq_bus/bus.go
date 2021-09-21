package rabbitmqbus

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	eventbus "github.com/paoloposso/crosstownbus/event_bus"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Uri string
}

type EventBus struct {
	channel *amqp.Channel
}

func CreateEventBus(config RabbitMQConfig) (eventbus.EventBus, error) {
	conn, err := amqp.Dial(config.Uri)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return EventBus{
		channel: channel,
	}, nil
}

func (pub EventBus) Publish(message interface{}) error {
	tp := reflect.TypeOf(message)
	eventName := tp.Name()
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

func (bus EventBus) Subscribe(event reflect.Type, eventHandler eventbus.EventHandler, resilienceOptions *eventbus.ResilienceOptions) error {
	eventName := event.Name()
	ch := bus.channel
	queueName := fmt.Sprintf("%s_%s_queue", eventName, reflect.TypeOf(eventHandler).Name())

	err := ch.ExchangeDeclare(eventName, "fanout", false, false, false, false, nil)
	if err != nil {
		return err
	}

	var retryOptions amqp.Table
	retryExchange := ""

	if resilienceOptions.MaxRetryTimes > 0 {
		deadLetterExchange, deadLetterRoutingKey, err := bus.createDeadLetter(queueName, reflect.TypeOf(eventHandler).Name(), resilienceOptions.RetrySeconds)
		if err != nil {
			return err
		}
		retryOptions = amqp.Table{
			"x-dead-letter-exchange":    deadLetterExchange,
			"x-dead-letter-routing-key": deadLetterRoutingKey,
		}
		retryExchange = deadLetterExchange
	}

	queue, err := ch.QueueDeclare(queueName, false, true, false, false, retryOptions)
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
		resilienceOptions == nil,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if resilienceOptions == nil || resilienceOptions.MaxRetryTimes == 0 {
		go func() {
			for msg := range msgs {
				go eventHandler.Handle(msg.Body)
			}
		}()
	} else {
		handleWithRetry(msgs, eventHandler, resilienceOptions.MaxRetryTimes)
	}
	return nil
}

func handleWithRetry(msgs <-chan amqp.Delivery, eventHandler eventbus.EventHandler, maxRetry uint16) {
	go func() {
		for msg := range msgs {
			if err := eventHandler.Handle(msg.Body); err != nil {
				log.Printf("Error handling message: %s", err)
				var redeliveryCount uint16 = 0
				if msg.Headers["x-redelivered-count"] != nil {
					redeliveryCount = msg.Headers["x-redelivered-count"].(uint16)
				}
				if redeliveryCount >= maxRetry {
					msg.Ack(false)
				} else {
					msg.Headers["x-redelivered-count"] = redeliveryCount + 1
					msg.Reject(false)
				}
			} else {
				msg.Ack(false)
			}
		}
	}()
}

func (bus EventBus) createDeadLetter(queueName string, handler string, retrySeconds uint16) (string, string, error) {
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
		"x-message-ttl":             retrySeconds,
	}

	if _, err := bus.channel.QueueDeclare(queueNameDl, true, true, false, false, options); err != nil {
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
