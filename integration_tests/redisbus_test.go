package integration_tests

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

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
		_ = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedSendMailHandler{}, nil)
		err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedIncludeHandler{}, nil)
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
		time.Sleep(2 * time.Second)
		err = bus.Publish(eventsamples.UserCreated{Name: "tes324t"})
		if err != nil {
			log.Fatalf("Error: %q", err)
		}
	}
	c := make(chan bool, 1)
	fmt.Printf("%t", <-c)
}
