package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/paoloposso/crosstownbus"
	eventsamples "github.com/paoloposso/crosstownbus/event_samples"
)

func main() {
	TestConsumeRabbit()
}

// main function, only for testing purpose for now
func TestConsumeRabbit() {
	_ = godotenv.Load()

	bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedHandler{})
		bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedHandler2{})

		time.Sleep(2 * time.Second)

		bus.Publish(eventsamples.UserCreated{Name: "test", Id: 55})
		bus.Publish(eventsamples.UserCreated{Name: "test", Id: 11243455})
		bus.Publish(eventsamples.UserCreated{Name: "test", Id: 5125})
		bus.Publish(eventsamples.UserCreated{Name: "test", Id: 55})
		bus.Publish(eventsamples.UserCreated{Name: "test", Id: 522325})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}
