package core

type EventHandler interface {
	Handle(event []byte) error
}
