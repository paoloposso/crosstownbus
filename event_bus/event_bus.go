package eventbus

import "reflect"

type EventBus interface {
	Publish(message interface{}) error
	Subscribe(event reflect.Type, eventHandler IntegrationEventHandler)
}
