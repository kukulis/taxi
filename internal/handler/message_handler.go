package handler

import (
	"fmt"

	"darbelis.eu/taxi/pkg/ws"
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

/**

Create a structure in the 'internal/handler' named WebSocketHubHandler, which contains  a map[MessageType]MessageHandler,

  where MessageHandler is func(clientId string, payload MessageInterface)

  It also has a method:
  func (hh*WebSocketHubHandler)  Handle(messageChannel <-chan ws.ClientMessage)

*/
