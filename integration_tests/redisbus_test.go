package integration_tests

import (
	"log"
	"reflect"
	"testing"

	"github.com/paoloposso/crosstownbus"
	eventsamples "github.com/paoloposso/crosstownbus/integration_tests/event_samples"
)

func TestCreateRedisEventBus(t *testing.T) {
	_, err := crosstownbus.CreateRedisEventBus("localhost:6379", "")
	if err != nil {
		t.Fatalf(`Error: %q`, err)
	}
}

func TestRedisPubSub(t *testing.T) {
	bus, err := crosstownbus.CreateRedisEventBus("localhost:6379", "")
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
