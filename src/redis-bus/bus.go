package redisbus

import (
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
	_ = bus.redisClient.Publish(bus.channel, message)
	return nil
}

func (bus Bus) SubscribeEvent(config interface{}, eventHandler eventbus.IntegrationEventHandler) {
	fmt.Println("started redis consume")
	go func() {
		for msg := range bus.redisClient.Subscribe(bus.channel).Channel() {
			fmt.Println(msg.Payload)
		}
		fmt.Println("exited redis channel")
	}()
}
