package rabbitmqbus

import (
	"encoding/json"
	"fmt"
	"log"

	eventbus "github.com/paoloposso/crosstownbus/event_bus"
	"github.com/streadway/amqp"
)

type RabbitMQOptions struct {
	Uri          string
	QueueName    string
	ExchangeName string
	// RoutingKey   string
	// ExchangeType string - only fanout exchanges for now. Other types will be added in the next version
}

type Bus struct {
	channel *amqp.Channel
	event   string
}

func CreateBus(eventName string, config RabbitMQOptions) (eventbus.Bus, error) {
	conn, err := amqp.Dial(config.Uri)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.ExchangeDeclare(eventName, "fanout", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return Bus{
		channel: channel,
		event:   eventName,
	}, nil
}

func (bus Bus) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ch := bus.channel
	err = ch.Publish(
		bus.event, // exchange
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

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	log.Println("Started RabbitMQ consume")
	ch := bus.channel
	go func() {
		queueName := fmt.Sprintf("%s_queue", bus.event)

		err := ch.ExchangeDeclare(bus.event, "fanout", false, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		err = ch.QueueBind(queue.Name, "", bus.event, false, nil)
		if err != nil {
			log.Fatal(err)
		}

		msgs, err := ch.Consume(
			queueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
		for msg := range msgs {
			go eventHandler.Handle(msg.Body)
		}
	}()
}
