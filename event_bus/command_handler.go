package eventbus

type CommandHandler interface {
	Handle(command []byte)
}
