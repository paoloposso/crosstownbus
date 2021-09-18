package redisbus

import (
	"log"

	eventbus "github.com/paoloposso/crosstown-bus/src/event_bus"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RabbitMQConfig struct {
	Uri      string
	Password string
}

type Bus struct {
	connection *amqp.Connection
	channel    string
}

func CreateBus(channel string, config RabbitMQConfig) (eventbus.Bus, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	return Bus{
		connection: conn,
		channel:    channel,
	}, nil
}

func (bus Bus) Publish(message interface{}) error {
	// str, err := json.Marshal(message)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// cmd := bus.redisClient.Publish(bus.channel, str)
	// if cmd.Err() != nil {
	// 	panic(cmd.Err().Error())
	// }
	return nil
}

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	// fmt.Println("started redis consume")
	// go func() {
	// 	cmd := bus.redisClient.Ping()
	// 	if cmd.Err() != nil {
	// 		panic(cmd.Err().Error())
	// 	}
	// 	for msg := range bus.redisClient.Subscribe(bus.channel).Channel() {
	// 		eventHandler.Handle([]byte(msg.Payload))
	// 	}
	// }()
}
