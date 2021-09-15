package eventbus

type IntegrationEvent interface {
	GetEventId() string
}
