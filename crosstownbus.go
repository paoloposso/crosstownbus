package crosstownbus

import (
	"github.com/paoloposso/crosstownbus/core"
	rabbitmqbus "github.com/paoloposso/crosstownbus/rabbitmq_bus"
	redisbus "github.com/paoloposso/crosstownbus/redis_bus"
)

func CreateRedisEventBus(uri string, password string) (core.EventBus, error) {
	return redisbus.CreateBus(
		redisbus.RedisConfig{
			Uri:      uri,
			Password: password,
		})
}

func CreateRabbitMQEventBus(uri string) (core.EventBus, error) {
	return rabbitmqbus.CreateEventBus(
		rabbitmqbus.RabbitMQConfig{
			Uri: uri,
		})
}
