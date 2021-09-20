package crosstownbus

import (
	eventbus "github.com/paoloposso/crosstownbus/event_bus"
	rabbitmqbus "github.com/paoloposso/crosstownbus/rabbitmq_bus"
	redisbus "github.com/paoloposso/crosstownbus/redis_bus"
)

func CreateRedisEventBus(uri string, password string) (eventbus.EventBus, error) {
	return redisbus.CreateBus(
		redisbus.RedisConfig{
			Uri:      uri,
			Password: password,
		})
}

func CreateRabbitMQEventBus(uri string) (eventbus.EventBus, error) {
	return rabbitmqbus.CreateEventBus(
		rabbitmqbus.RabbitMQConfig{
			Uri: uri,
		})
}
