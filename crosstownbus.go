package crosstownbus

import (
	"reflect"

	eventbus "github.com/paoloposso/crosstownbus/event_bus"
	rabbitmqbus "github.com/paoloposso/crosstownbus/rabbitmq_bus"
	redisbus "github.com/paoloposso/crosstownbus/redis_bus"
)

func CreateRedisEventBus(event reflect.Type, uri string, password string) (eventbus.EventBus, error) {
	return redisbus.CreateBus(
		event.Name(),
		redisbus.RedisConfig{
			Uri:      uri,
			Password: password,
		})
}

func CreateRabbitMQEventBus(uri string) (eventbus.EventBus, error) {
	return rabbitmqbus.CreateBus(
		rabbitmqbus.RabbitMQConfig{
			Uri: uri,
		})
}
