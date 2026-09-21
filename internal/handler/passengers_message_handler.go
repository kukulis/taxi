package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/ws"
)

type PassengersMessageHandler struct {
}

func (p PassengersMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel
		fmt.Println("PassengersMessageHandler:", message)
	}
}

// forcing to implement interface
var _ MessageHandler = &PassengersMessageHandler{}

func NewPassengersMessageHandler() *PassengersMessageHandler {
	return &PassengersMessageHandler{}
}
