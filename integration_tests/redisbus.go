package main

func TestRedisPubSub() {
	// bus, err := crosstownbus.CreateRedisEventBus("localhost:6379", "")
	// if err != nil {
	// 	log.Fatalf("Error: %q", err)
	// } else {
	// 	_ = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedSendMailHandler{}, nil)
	// 	_ = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedIncludeHandler{}, nil)
	// 	err = bus.Subscribe(reflect.TypeOf(eventsamples.UserCreated{}), eventsamples.UserCreatedIncludeHandler{}, nil)
	// 	if err != nil {
	// 		log.Fatalf("Error: %q", err)
	// 	}
	// 	time.Sleep(2 * time.Second)
	// 	_ = bus.Publish(eventsamples.UserCreated{Name: "tes324t"})
	// 	_ = bus.Publish(eventsamples.UserCreated{Name: "xHH-90"})
	// 	_ = bus.Publish(eventsamples.UserCreated{Name: "paolo"})
	// 	err = bus.Publish(eventsamples.UserCreated{Name: "test3"})
	// 	if err != nil {
	// 		log.Fatalf("Error: %q", err)
	// 	}
	// }
}
