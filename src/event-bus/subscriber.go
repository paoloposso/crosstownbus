package eventbus

type Subscriber interface {
	SubscribeEvent(config struct{}, eventName string, eventHandler IntegrationEventHandler)
}
