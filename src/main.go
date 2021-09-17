package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	redisbus "github.com/paoloposso/crosstown-bus/src/redis-bus"
)

func main() {
	_ = godotenv.Load()

	bus := redisbus.CreateBus("test1")
	func() {
		bus.Subscribe(HandlerSample{})
	}()

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
