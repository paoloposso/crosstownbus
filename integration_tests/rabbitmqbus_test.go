package integration_tests

import (
	"log"
	"reflect"
	"testing"

	"github.com/paoloposso/crosstownbus"
	eventsamples "github.com/paoloposso/crosstownbus/integration_tests/event_samples"
)

func TestCreateRabbitEventBus(t *testing.T) {
	_, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf(`Error: %q`, err)
	}
}

func TestRabbitPubSub(t *testing.T) {
	bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("Error: %q", err)
	} else {
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedHandler{})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
		err = bus.Publish(eventsamples.UserCreated{Name: "tes324t"})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
	}
}
