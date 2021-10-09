package core

import "reflect"

type CommandBus interface {
	Send(message interface{}) error
	Receive(command reflect.Type, eventHandler CommandHandler) error
}
