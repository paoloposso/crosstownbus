package crosstownbus

import (
	"github.com/paoloposso/crosstownbus/core"
	rabbitmqbus "github.com/paoloposso/crosstownbus/rabbitmq_bus"
)

func CreateRabbitMQEventBus(uri string) (core.EventBus, error) {
	return rabbitmqbus.CreateEventBus(
		rabbitmqbus.RabbitMQConfig{
			Uri: uri,
		})
}
