package redis

import eventbus "github.com/paoloposso/crosstown-bus/src/event-bus"

type Bus struct{}

func CreateBus() eventbus.Bus {
	return Bus{}
}

func (bus Bus) Publish(integrationEvent eventbus.IntegrationEvent) error {
	return nil
}

func (bus Bus) SubscribeEvent(config struct{}, eventHandler eventbus.IntegrationEventHandler) {}
