package eventbus

type Bus interface {
	Publish(integrationEvent IntegrationEvent) error
	SubscribeEvent(config struct{}, eventHandler IntegrationEventHandler)
}
