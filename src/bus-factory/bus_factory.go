package busfactory

import (
	"reflect"

	eventbus "github.com/pvictorsys/crosstown-bus/src/event_bus"
	redisbus "github.com/pvictorsys/crosstown-bus/src/redis_bus"
)

func CreateRedisBus(event reflect.Type, uri string, password string) (eventbus.Bus, error) {
	config := redisbus.RedisConfig{Uri: uri, Password: password}
	return redisbus.CreateBus(event.Name(), config)
}
