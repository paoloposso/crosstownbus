package core

type CommandHandler interface {
	Handle(command []byte)
}
