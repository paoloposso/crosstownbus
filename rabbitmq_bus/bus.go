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
	//event   string
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
		log.Fatal(err)
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

func (bus EventBus) Subscribe(event reflect.Type, eventHandler eventbus.IntegrationEventHandler) {
	log.Println("Started RabbitMQ consume")
	eventName := event.Name()
	ch := bus.channel
	go func() {
		queueName := fmt.Sprintf("%s_%s_queue", eventName, reflect.TypeOf(eventHandler).Name())

		err := ch.ExchangeDeclare(eventName, "fanout", false, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		err = ch.QueueBind(queue.Name, "", eventName, false, nil)
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
