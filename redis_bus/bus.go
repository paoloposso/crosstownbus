package redisbus

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	eventbus "github.com/paoloposso/crosstown-bus/event_bus"
)

type RedisConfig struct {
	Uri      string
	Password string
}

type Bus struct {
	redisClient *redis.Client
	channel     string
}

func CreateBus(channel string, config RedisConfig) (eventbus.Bus, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Uri,
		Password: config.Password,
		DB:       0,
	})
	cmd := client.Ping()
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return Bus{
		redisClient: client,
		channel:     channel,
	}, nil
}

func (bus Bus) Publish(message interface{}) error {
	str, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd := bus.redisClient.Publish(bus.channel, str)
	if cmd.Err() != nil {
		panic(cmd.Err().Error())
	}
	return nil
}

func (bus Bus) Subscribe(eventHandler eventbus.IntegrationEventHandler) {
	fmt.Println("started redis consume")
	go func() {
		cmd := bus.redisClient.Ping()
		if cmd.Err() != nil {
			panic(cmd.Err().Error())
		}
		for msg := range bus.redisClient.Subscribe(bus.channel).Channel() {
			go eventHandler.Handle([]byte(msg.Payload))
		}
	}()
}
