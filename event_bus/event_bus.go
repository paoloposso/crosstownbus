package eventbus

import "reflect"

type EventBus interface {
	Publish(event reflect.Type, message interface{}) error
	Subscribe(event reflect.Type, eventHandler IntegrationEventHandler)
}
