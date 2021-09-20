package eventbus

import "reflect"

type EventBus interface {
	Publish(message interface{}) error
	Subscribe(event reflect.Type, handler IntegrationEventHandler) error
}
