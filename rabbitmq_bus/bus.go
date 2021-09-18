package rabbitmqbus

import (
	"encoding/json"
	"fmt"
	"log"

	eventbus "github.com/paoloposso/crosstown-bus/event_bus"
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
	connection   *amqp.Connection
	channel      string
	exchangeName string
}

func CreateBus(channel string, config RabbitMQOptions) (eventbus.Bus, error) {
	conn, err := amqp.Dial(config.Uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	return Bus{
		connection:   conn,
		channel:      channel,
		exchangeName: config.ExchangeName,
	}, nil
}

func (bus Bus) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ch, err := bus.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	ch.ExchangeDeclare(bus.exchangeName, "fanout", false, false, false, false, nil)

	err = ch.Publish(
		bus.exchangeName, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message to RabbitMQ")
	return nil
}

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	log.Println("Started RabbitMQ consume")
	go func() {
		ch, err := bus.connection.Channel()
		failOnError(err, "Failed to open a channel")
		queueName := fmt.Sprintf("%s_queue", bus.channel)
		ch.QueueDeclare(queueName, false, false, false, false, nil)
		ch.QueueBind(queueName, "", bus.exchangeName, false, nil)
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
			eventHandler.Handle(msg.Body)
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
