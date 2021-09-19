package integration_tests

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

// main function, only for testing purpose for now
func TestConsumeRedis() {
	_ = godotenv.Load()

	bus, err := crosstownbus.CreateRedisEventBus("localhost:6379", "")

	errs := make(chan error, 1)

	if err != nil {
		errs <- err
	} else {
		bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedHandler{})
		bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedHandler2{})
		time.Sleep(2 * time.Second)
		bus.Publish(eventsamples.UserCreated{Name: "tes324t"})
		bus.Publish(eventsamples.UserCreated{Name: "asddsdsad"})
		bus.Publish(eventsamples.UserCreated{Name: "12TOEf"})
		bus.Publish(eventsamples.UserCreated{Name: "ZZZZZZZ"})
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Terminated. Reason:", <-errs)
}
