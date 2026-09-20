package ws

import "fmt"

type MessageHandler interface {
	Handle(message ClientMessage)
}

type SimpleMessageHandler struct {
}

func (handler *SimpleMessageHandler) Handle(message ClientMessage) {
	fmt.Println("Simple message handler:", message)
}
