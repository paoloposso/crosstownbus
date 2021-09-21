package integration_tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/paoloposso/crosstownbus"
	eventbus "github.com/paoloposso/crosstownbus/event_bus"
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
		t.Fatalf("Error: %q", err)
	} else {
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedSendMailHandler{}, &eventbus.ResilienceOptions{RetrySeconds: 5, MaxRetryTimes: 3})
		if err != nil {
			t.Fatalf("Error: %q", err)
		}
		time.Sleep(2 * time.Second)
		err = bus.Publish(eventsamples.UserCreated{Name: "error"})
		_ = bus.Publish(eventsamples.UserCreated{Name: "ok"})
		if err != nil {
			t.Fatalf("Error: %q", err)
		}
	}
	ch := make(chan bool, 1)
	fmt.Printf("%t", <-ch)
}
