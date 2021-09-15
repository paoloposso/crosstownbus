package eventbus

type Publisher interface {
	Publish(integrationEvent IntegrationEvent) error
}
