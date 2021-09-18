package main

import (
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

	bus, _ := busfactory.CreateRedisBus(reflect.TypeOf(UserCreated{}), "localhost:6379", "")

	bus.Subscribe(HandlerSample{})
	bus.Publish(UserCreated{Name: "test"})

	errs := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated ", <-errs)
}

type UserCreated struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

type HandlerSample struct{}

func (handler HandlerSample) Handle(event string) {
	fmt.Println(event)
}
