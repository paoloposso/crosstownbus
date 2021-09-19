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
	"github.com/paoloposso/crosstownbus"
)

type UserCreated struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

type UserCreatedHandler struct{}

func (handler UserCreatedHandler) Handle(event []byte) {
	var user *UserCreated
	json.Unmarshal(event, &user)
	fmt.Println(user.Name, "received:", time.Now())
	time.Sleep(5 * time.Second)
}

type UserCreatedHandler2 struct{}

func (handler UserCreatedHandler2) Handle(event []byte) {
	var user *UserCreated
	json.Unmarshal(event, &user)
	fmt.Println(user.Name, "received:", time.Now())
	time.Sleep(5 * time.Second)
}

// main function, only for testing purpose for now
func main() {
	_ = godotenv.Load()

	bus, err := crosstownbus.CreateRabbitMQEventBus("amqp://guest:guest@localhost:5672/")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(reflect.TypeOf(UserCreated{}), UserCreatedHandler{})
		bus.Subscribe(reflect.TypeOf(UserCreated{}), UserCreatedHandler2{})

		time.Sleep(2 * time.Second)

		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "tes324t"})
		bus.Publish(reflect.TypeOf(UserCreated{}), UserCreated{Name: "test", Id: 55})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}
