package crosstownbus

import (
	"reflect"

	eventbus "github.com/paoloposso/crosstownbus/event_bus"
	redisbus "github.com/paoloposso/crosstownbus/redis_bus"
)

func CreateRedisBus(event reflect.Type, uri string, password string) (eventbus.Bus, error) {
	config := redisbus.RedisConfig{Uri: uri, Password: password}
	return redisbus.CreateBus(event.Name(), config)
}
