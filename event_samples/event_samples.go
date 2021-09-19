package eventsamples

import (
	"encoding/json"
	"fmt"
	"time"
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
	fmt.Println(user.Name, "handled 2:", time.Now())
	time.Sleep(5 * time.Second)
}
