package eventbus

type Bus interface {
	Publish(message interface{}) error
	Subscribe(eventHandler IntegrationEventHandler)
}
