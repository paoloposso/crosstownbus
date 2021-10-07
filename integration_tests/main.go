package main

import (
	"fmt"
)

func main() {
	TestRabbitPubSub()
	TestRedisPubSub()
	ch := make(chan bool, 3)
	fmt.Printf("%t", <-ch)
}
