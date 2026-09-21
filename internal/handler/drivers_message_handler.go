package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/ws"
)

type DriversMessageHandler struct {
}

func (d DriversMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel
		fmt.Println("DriversMessageHandler:", message)
	}
}

// forcing to implement interface
var _ MessageHandler = &DriversMessageHandler{}

func NewDriversMessageHandler() *DriversMessageHandler {
	return &DriversMessageHandler{}
}
