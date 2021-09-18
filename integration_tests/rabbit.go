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
)

// main function, only for testing purpose for now
func main() {
	_ = godotenv.Load()

	bus, err := crosstownbus.CreateRabbitMQBus(reflect.TypeOf(UserCreated{}), "amqp://guest:guest@localhost:5672/")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(HandlerSample{})

		time.Sleep(2 * time.Second)

		bus.Publish(UserCreated{Name: "tes324t"})
		bus.Publish(UserCreated{Name: "test", Id: 55})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}
