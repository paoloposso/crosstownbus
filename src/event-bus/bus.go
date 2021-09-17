package eventbus

type Bus interface {
	Publish(message interface{}) error
	SubscribeEvent(config interface{}, eventHandler IntegrationEventHandler)
}
