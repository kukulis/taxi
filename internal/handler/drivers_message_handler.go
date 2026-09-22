package handler

import (
	"fmt"
	"log/slog"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
)

type DriversMessageHandler struct {
	mainState     *state.MainState
	driversHub    ws.HubInterface
	passengersHub ws.HubInterface
	logger        *slog.Logger
}

func NewDriversMessageHandler(mainState *state.MainState, driversHub ws.HubInterface, passengersHub ws.HubInterface) *DriversMessageHandler {
	return &DriversMessageHandler{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
		logger:        slog.Default().With("handler", "DriversMessageHandler"),
	}
}

func (d DriversMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {
	for {
		message := <-messageChannel

		switch msg := message.Message.(type) {
		case *messages.ClientRespondsCoordinates:
			d.handleClientRespondsCoordinates(message.ClientId, msg)
		case *messages.ClientRespondsInfo:
			d.handleClientRespondsInfo(message.ClientId, msg)
		case *messages.ClientResponseError:
			d.handleClientResponseError(message.ClientId, msg)
		case *messages.ClientError:
			d.handleClientError(message.ClientId, msg)
		case *messages.DriverCancelsOffer:
			d.handleDriverCancelsOffer(message.ClientId, msg)
		case *messages.DriverAcceptsOffer:
			d.handleDriverAcceptsOffer(message.ClientId, msg)
		case *messages.DriverRejectsOffer:
			d.handleDriverRejectsOffer(message.ClientId, msg)
		case *messages.DriverChangesStatus:
			d.handleDriverChangesStatus(message.ClientId, msg)
		default:
			d.logger.Warn("unexpected message type",
				"client_id", message.ClientId,
				"message_type", fmt.Sprintf("%T", message.Message),
			)
		}
	}
}

func (d DriversMessageHandler) handleClientRespondsCoordinates(clientId string, msg *messages.ClientRespondsCoordinates) {
	found := d.mainState.UpdateDriver(clientId, func(driver *state.Driver) {
		driver.Lat = msg.Lat
		driver.Lon = msg.Lon
		driver.CoordinatesReceivedAt = d.mainState.Clock.Now()
	})
	if !found {
		d.logger.Warn("ClientRespondsCoordinates from unknown driver", "client_id", clientId)
	}
}

func (d DriversMessageHandler) handleClientRespondsInfo(clientId string, msg *messages.ClientRespondsInfo) {
	found := d.mainState.UpdateDriver(clientId, func(driver *state.Driver) {
		driver.Phone = msg.Phone
		driver.VehicleInfo = msg.VehicleInfo
		if msg.Ip != "" {
			driver.Ip = msg.Ip
		}
		driver.InfoReceivedAt = d.mainState.Clock.Now()
	})
	if !found {
		d.logger.Warn("ClientRespondsInfo from unknown driver", "client_id", clientId)
	}
}

// handleClientResponseError logs the driver client's report that it didn't
// understand a request the server sent it. No MainState field maps to this,
// so it's just logged for now.
func (d DriversMessageHandler) handleClientResponseError(clientId string, msg *messages.ClientResponseError) {
	d.logger.Warn("client didn't understand a server request", "client_id", clientId, "error", msg.Error)
}

// handleClientError logs the driver client's own functionality error (e.g. it
// can't read device coordinates). No MainState field maps to this, so it's
// just logged for now.
func (d DriversMessageHandler) handleClientError(clientId string, msg *messages.ClientError) {
	d.logger.Error("client reported a functionality error", "client_id", clientId, "error", msg.Error)
}

func (d DriversMessageHandler) handleDriverCancelsOffer(clientId string, msg *messages.DriverCancelsOffer) {
	d.logger.Info("driver cancels offer", "client_id", clientId, "passenger_id", msg.PassengerId)
}

func (d DriversMessageHandler) handleDriverAcceptsOffer(clientId string, msg *messages.DriverAcceptsOffer) {
	d.logger.Info("driver accepts offer", "client_id", clientId, "passenger_id", msg.PassengerId)
}

func (d DriversMessageHandler) handleDriverRejectsOffer(clientId string, msg *messages.DriverRejectsOffer) {
	d.logger.Info("driver rejects offer", "client_id", clientId, "passenger_id", msg.PassengerId)
}

func (d DriversMessageHandler) handleDriverChangesStatus(clientId string, msg *messages.DriverChangesStatus) {
	d.logger.Info("driver changes status", "client_id", clientId, "status", msg.Status)
}

// forcing to implement interface
var _ MessageHandler = &DriversMessageHandler{}
