package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/joho/godotenv"
	busfactory "github.com/paoloposso/crosstown-bus/src/bus-factory"
)

func main() {
	_ = godotenv.Load()

	bus, err := busfactory.CreateRedisBus(reflect.TypeOf(UserCreated{}), "localhostx:6379", "")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(HandlerSample{})

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
	fmt.Println(user.Name)
}
