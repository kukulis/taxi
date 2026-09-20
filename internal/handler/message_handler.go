package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/ws"
)

type MessageHandler interface {
	Handle(messageChannel <-chan ws.ClientMessage)
}

type SimpleMessageHandler struct {
}

func (handler *SimpleMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel
		fmt.Println("Simple message handler:", message)
	}
}
