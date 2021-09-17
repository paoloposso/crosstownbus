package redisbus

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	eventbus "github.com/paoloposso/crosstown-bus/src/event-bus"
)

type Bus struct {
	redisClient *redis.Client
	channel     string
}

func CreateBus(channel string) eventbus.Bus {
	return Bus{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
		channel: channel,
	}
}

func (bus Bus) Publish(message interface{}) error {
	str, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	_ = bus.redisClient.Publish(bus.channel, str)
	return nil
}

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	fmt.Println("started redis consume")
	go func() {
		for msg := range bus.redisClient.Subscribe(bus.channel).Channel() {
			eventHandler.Handle(msg.Payload)
		}
		fmt.Println("exited redis channel")
	}()
}
