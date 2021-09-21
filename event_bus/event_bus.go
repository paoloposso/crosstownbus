package eventbus

import "reflect"

// For now, Retry is only supported when using RabbitMQ
type ResilienceOptions struct {
	RetrySeconds  uint16 // Seconds before retrying a failed handling
	MaxRetryTimes uint16 // Set zero for not determining a number of retries
}

type EventBus interface {
	Publish(message interface{}) error
	Subscribe(event reflect.Type, handler EventHandler, resilienceOptions *ResilienceOptions) error
}
