package eventsamples

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type UserCreated struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

type UserCreatedSendMailHandler struct{}

func (handler UserCreatedSendMailHandler) Handle(event []byte) error {
	var user *UserCreated
	err := json.Unmarshal(event, &user)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	if user.Name == "error" {
		return errors.New("error")
	}

	fmt.Println(user.Name, "received and email sent at:", time.Now())

	return nil
}

type UserCreatedIncludeHandler struct{}

func (handler UserCreatedIncludeHandler) Handle(event []byte) error {
	var user *UserCreated
	err := json.Unmarshal(event, &user)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	if user.Name == "error" {
		return errors.New("error")
	}

	fmt.Println(user.Name, "received and included:", time.Now())

	return nil
}
