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
func mainx() {
	_ = godotenv.Load()

	bus, err := crosstownbus.CreateRedisEventBus(reflect.TypeOf(UserCreated{}), "localhost:6379", "")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(reflect.TypeOf(UserCreated{}), UserCreatedHandler{})

		time.Sleep(2 * time.Second)

		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "tes324t"})
		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "test"})
		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "tes34t"})
		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "test4"})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}
