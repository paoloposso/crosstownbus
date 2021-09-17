package main

import (
	"fmt"

	"github.com/joho/godotenv"
	redisbus "github.com/paoloposso/crosstown-bus/src/redis-bus"
)

func main() {
	_ = godotenv.Load()

	bus := redisbus.CreateBus("test1")
	func() {
		bus.SubscribeEvent(nil, HandlerSample{})
	}()

	bus.Publish("aaaaaa")
	bus.Publish("aaaa3343aa")
	bus.Publish("aaaaaa")
	bus.Publish("aaaaaa")
	bus.Publish("232343")
	bus.Publish("aaaaaa")
	bus.Publish("aaaaf5434aa")
	bus.Publish("aaaaaa")

	<-make(chan int)
}

type HandlerSample struct{}

func (handler HandlerSample) Handle(event struct{}) {
	fmt.Println(event)
}
