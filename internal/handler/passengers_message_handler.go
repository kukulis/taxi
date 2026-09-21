package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
)

type PassengersMessageHandler struct {
	mainState *state.MainState
}

func NewPassengersMessageHandler(mainState *state.MainState) *PassengersMessageHandler {
	return &PassengersMessageHandler{mainState: mainState}
}

func (p PassengersMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel

		switch msg := message.Message.(type) {
		case messages.ClientRespondsCoordinates:
			p.handleClientRespondsCoordinates(message.ClientId, msg)
		case messages.PassengerInvitesDriver:
			p.handlePassengerInvitesDriver(message.ClientId, msg)
		case messages.PassengerCancelsInvite:
			p.handlePassengerCancelsInvite(message.ClientId, msg)
		case messages.PassengerRequestDriverCoords:
			p.handlePassengerRequestDriverCoords(message.ClientId, msg)
		default:
			fmt.Printf("PassengersMessageHandler: unexpected message type %T from client %s, skipping\n", message.Message, message.ClientId)
		}
	}
}

func (p PassengersMessageHandler) handleClientRespondsCoordinates(clientId string, msg messages.ClientRespondsCoordinates) {
	fmt.Println("PassengersMessageHandler: ClientRespondsCoordinates from", clientId, msg)
}

func (p PassengersMessageHandler) handlePassengerInvitesDriver(clientId string, msg messages.PassengerInvitesDriver) {
	fmt.Println("PassengersMessageHandler: PassengerInvitesDriver from", clientId, msg)
}

func (p PassengersMessageHandler) handlePassengerCancelsInvite(clientId string, msg messages.PassengerCancelsInvite) {
	fmt.Println("PassengersMessageHandler: PassengerCancelsInvite from", clientId, msg)
}

func (p PassengersMessageHandler) handlePassengerRequestDriverCoords(clientId string, msg messages.PassengerRequestDriverCoords) {
	fmt.Println("PassengersMessageHandler: PassengerRequestDriverCoords from", clientId, msg)
}

// forcing to implement interface
var _ MessageHandler = &PassengersMessageHandler{}
