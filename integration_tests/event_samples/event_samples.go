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

type UserCreatedSendMailHandler struct{}

func (handler UserCreatedSendMailHandler) Handle(event []byte) {
	var user *UserCreated
	json.Unmarshal(event, &user)
	fmt.Println(user.Name, "received:", time.Now())
	time.Sleep(5 * time.Second)
}
