package busfactory

import (
	"reflect"

	eventbus "github.com/paoloposso/crosstown-bus/event_bus"
	redisbus "github.com/paoloposso/crosstown-bus/redis_bus"
)

func CreateRedisBus(event reflect.Type, uri string, password string) (eventbus.Bus, error) {
	config := redisbus.RedisConfig{Uri: uri, Password: password}
	return redisbus.CreateBus(event.Name(), config)
}
