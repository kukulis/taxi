package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
)

type DriversMessageHandler struct {
	mainState     *state.MainState
	driversHub    ws.HubInterface
	passengersHub ws.HubInterface
}

func NewDriversMessageHandler(mainState *state.MainState, driversHub ws.HubInterface, passengersHub ws.HubInterface) *DriversMessageHandler {
	return &DriversMessageHandler{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
	}
}

func (d DriversMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel

		switch msg := message.Message.(type) {
		case messages.ClientRespondsCoordinates:
			d.handleClientRespondsCoordinates(message.ClientId, msg)
		case messages.ClientRespondsInfo:
			d.handleClientRespondsInfo(message.ClientId, msg)
		case messages.DriverCancelsOffer:
			d.handleDriverCancelsOffer(message.ClientId, msg)
		case messages.DriverAcceptsOffer:
			d.handleDriverAcceptsOffer(message.ClientId, msg)
		case messages.DriverRejectsOffer:
			d.handleDriverRejectsOffer(message.ClientId, msg)
		case messages.DriverChangesStatus:
			d.handleDriverChangesStatus(message.ClientId, msg)
		default:
			fmt.Printf("DriversMessageHandler: unexpected message type %T from client %s, skipping\n", message.Message, message.ClientId)
		}
	}
}

func (d DriversMessageHandler) handleClientRespondsCoordinates(clientId string, msg messages.ClientRespondsCoordinates) {
	found := d.mainState.UpdateDriver(clientId, func(driver *state.Driver) {
		driver.Lat = msg.Lat
		driver.Lon = msg.Lon
	})
	if !found {
		fmt.Println("DriversMessageHandler: ClientRespondsCoordinates from unknown driver", clientId)
	}
}

func (d DriversMessageHandler) handleClientRespondsInfo(clientId string, msg messages.ClientRespondsInfo) {
	found := d.mainState.UpdateDriver(clientId, func(driver *state.Driver) {
		driver.Phone = msg.Phone
		driver.VehicleInfo = msg.VehicleInfo
	})
	if !found {
		fmt.Println("DriversMessageHandler: ClientRespondsInfo from unknown driver", clientId)
	}
}

func (d DriversMessageHandler) handleDriverCancelsOffer(clientId string, msg messages.DriverCancelsOffer) {
	fmt.Println("DriversMessageHandler: DriverCancelsOffer from", clientId, msg)
}

func (d DriversMessageHandler) handleDriverAcceptsOffer(clientId string, msg messages.DriverAcceptsOffer) {
	fmt.Println("DriversMessageHandler: DriverAcceptsOffer from", clientId, msg)
}

func (d DriversMessageHandler) handleDriverRejectsOffer(clientId string, msg messages.DriverRejectsOffer) {
	fmt.Println("DriversMessageHandler: DriverRejectsOffer from", clientId, msg)
}

func (d DriversMessageHandler) handleDriverChangesStatus(clientId string, msg messages.DriverChangesStatus) {
	fmt.Println("DriversMessageHandler: DriverChangesStatus from", clientId, msg)
}

// forcing to implement interface
var _ MessageHandler = &DriversMessageHandler{}
