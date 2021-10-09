package redisbus

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-redis/redis"
	"github.com/paoloposso/crosstownbus/core"
)

type RedisConfig struct {
	Uri      string
	Password string
}

type EventBus struct {
	redisClient *redis.Client
}

func CreateBus(config RedisConfig) (core.EventBus, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Uri,
		Password: config.Password,
		DB:       0,
	})
	cmd := client.Ping()
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return EventBus{
		redisClient: client,
	}, nil
}

func (bus EventBus) Publish(message interface{}) error {
	tp := reflect.TypeOf(message)
	eventName := tp.Name()
	str, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd := bus.redisClient.Publish(eventName, str)
	if cmd.Err() != nil {
		return err
	}
	return nil
}

func (bus EventBus) Subscribe(event reflect.Type, eventHandler core.EventHandler, resilienceOptions *core.ResilienceOptions) error {
	cmd := bus.redisClient.Ping()
	if cmd.Err() != nil {
		return cmd.Err()
	}
	go func() {
		for msg := range bus.redisClient.Subscribe(event.Name()).Channel() {
			go eventHandler.Handle([]byte(msg.Payload))
		}
	}()
	return nil
}
