package eventbus

type EventHandler interface {
	Handle(event []byte) error
}
