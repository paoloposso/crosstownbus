package integration_tests

import (
	"testing"

	"github.com/paoloposso/crosstownbus"
)

func TestCreateRabbitEventBus(t *testing.T) {
	_, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf(`Error: %q`, err)
	}
}

func TestCreateRedisEventBus(t *testing.T) {
	_, err := crosstownbus.CreateRedisEventBus("localhost:6379", "")
	if err != nil {
		t.Fatalf(`Error: %q`, err)
	}
}
