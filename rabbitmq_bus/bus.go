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
	failOnError(err, "Failed to connect to RabbitMQ")
	channel, err := conn.Channel()
	failOnError(err, "Failed to connect to RabbitMQ")
	return Bus{
		channel: channel,
		event:   eventName,
	}, nil
}

func (bus Bus) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ch := bus.channel

	err = ch.ExchangeDeclare(bus.event, "fanout", false, false, false, false, nil)
	failOnError(err, "Failed creating exchange")

	err = ch.Publish(
		bus.event, // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message to RabbitMQ")
	return nil
}

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	log.Println("Started RabbitMQ consume")
	ch := bus.channel
	go func() {
		queueName := fmt.Sprintf("%s_queue", bus.event)

		err := ch.ExchangeDeclare(bus.event, "fanout", false, false, false, false, nil)
		failOnError(err, "Failed creating exchange")
		queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
		failOnError(err, "Failed creating queue")
		err = ch.QueueBind(queue.Name, "", bus.event, false, nil)
		failOnError(err, "Failed binding queue")

		msgs, err := ch.Consume(
			queueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to start consume")
		for msg := range msgs {
			go eventHandler.Handle(msg.Body)
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
