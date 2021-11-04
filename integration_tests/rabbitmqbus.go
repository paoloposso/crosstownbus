package main

import (
	"log"
	"reflect"
	"time"

	"github.com/paoloposso/crosstownbus"
	"github.com/paoloposso/crosstownbus/core"
	eventsamples "github.com/paoloposso/crosstownbus/integration_tests/event_samples"
)

func TestRabbitPubSub() {
	bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("Error: %q", err)
	} else {
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedSendMailHandler{}, &core.RetryOptions{RetrySeconds: 5, MaxRetryTimes: 3})
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedIncludeHandler{}, &core.RetryOptions{RetrySeconds: 5, MaxRetryTimes: 3})
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedIncludeHandler{}, &core.RetryOptions{RetrySeconds: 5, MaxRetryTimes: 3})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
		time.Sleep(2 * time.Second)
		err = bus.Publish(eventsamples.UserCreated{Name: "error"})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
		err = bus.Publish(eventsamples.UserCreated{Name: "Paolo"})
		// _ = bus.Publish(eventsamples.UserCreated{Name: "ok"})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
	}
}
