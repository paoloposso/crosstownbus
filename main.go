package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	busfactory "github.com/pvictorsys/crosstown-bus/src/bus-factory"
)

// main function, only for testing purpose for now
func main() {
	_ = godotenv.Load()

	bus, err := busfactory.CreateRedisBus(reflect.TypeOf(UserCreated{}), "localhost:6379", "")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(HandlerSample{})

		time.Sleep(2 * time.Second)

		bus.Publish(UserCreated{Name: "tes324t"})
		bus.Publish(UserCreated{Name: "test"})
		bus.Publish(UserCreated{Name: "tes34t"})
		bus.Publish(UserCreated{Name: "test4"})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}

type UserCreated struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

type HandlerSample struct{}

func (handler HandlerSample) Handle(event []byte) {
	var user *UserCreated
	json.Unmarshal(event, &user)
	fmt.Println(user.Name, "received:", time.Now())
	time.Sleep(5 * time.Second)
}
