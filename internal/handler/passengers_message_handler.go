package handler

import (
	"fmt"
	"log/slog"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

type PassengersMessageHandler struct {
	mainState     *state.MainState
	driversHub    ws2.HubInterface
	passengersHub ws2.HubInterface
	logger        *slog.Logger
}

func NewPassengersMessageHandler(mainState *state.MainState, driversHub ws2.HubInterface, passengersHub ws2.HubInterface) *PassengersMessageHandler {
	return &PassengersMessageHandler{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
		logger:        slog.Default().With("handler", "PassengersMessageHandler"),
	}
}

func (p *PassengersMessageHandler) Handle(messageChannel <-chan ws2.ClientMessage) {
	for {
		message := <-messageChannel

		switch msg := message.Message.(type) {
		case *messages.ClientRespondsCoordinates:
			p.handleClientRespondsCoordinates(message.ClientId, msg)
		case *messages.ClientRespondsInfo:
			p.handleClientRespondsInfo(message.ClientId, msg)
		case *messages.ClientResponseError:
			p.handleClientResponseError(message.ClientId, msg)
		case *messages.ClientError:
			p.handleClientError(message.ClientId, msg)
		case *messages.PassengerInvitesDriver:
			p.handlePassengerInvitesDriver(message.ClientId, msg)
		case *messages.PassengerCancelsInvite:
			p.handlePassengerCancelsInvite(message.ClientId, msg)
		case *messages.PassengerRequestDriverCoords:
			p.handlePassengerRequestDriverCoords(message.ClientId, msg)
		default:
			p.logger.Warn("unexpected message type",
				"client_id", message.ClientId,
				"message_type", fmt.Sprintf("%T", message.Message),
			)
		}
	}
}

func (p *PassengersMessageHandler) handleClientRespondsCoordinates(clientId string, msg *messages.ClientRespondsCoordinates) {
	found := p.mainState.UpdatePassenger(clientId, func(passenger *state.Passenger) {
		passenger.Lat = msg.Lat
		passenger.Lon = msg.Lon
		passenger.CoordinatesReceivedAt = p.mainState.Clock.Now()
	})
	if !found {
		p.logger.Warn("ClientRespondsCoordinates from unknown passenger", "client_id", clientId)
	}
}

func (p *PassengersMessageHandler) handleClientRespondsInfo(clientId string, msg *messages.ClientRespondsInfo) {
	found := p.mainState.UpdatePassenger(clientId, func(passenger *state.Passenger) {
		passenger.Phone = msg.Phone
		if msg.Ip != "" {
			passenger.Ip = msg.Ip
		}
		passenger.InfoReceivedAt = p.mainState.Clock.Now()
	})
	if !found {
		p.logger.Warn("ClientRespondsInfo from unknown passenger", "client_id", clientId)
	}
}

// handleClientResponseError logs the passenger client's report that it didn't
// understand a request the server sent it. No MainState field maps to this,
// so it's just logged for now.
func (p *PassengersMessageHandler) handleClientResponseError(clientId string, msg *messages.ClientResponseError) {
	p.logger.Warn("client didn't understand a server request", "client_id", clientId, "error", msg.Error)
}

// handleClientError logs the passenger client's own functionality error. No
// MainState field maps to this, so it's just logged for now.
func (p *PassengersMessageHandler) handleClientError(clientId string, msg *messages.ClientError) {
	p.logger.Error("client reported a functionality error", "client_id", clientId, "error", msg.Error)
}

func (p *PassengersMessageHandler) handlePassengerInvitesDriver(clientId string, msg *messages.PassengerInvitesDriver) {
	p.logger.Info("passenger invites driver", "client_id", clientId, "driver_id", msg.DriverId, "lat", msg.Lat, "lon", msg.Lon)

	invitationId, err := util.RandomHash(8)
	if err != nil {
		p.logger.Error("failed to create invitation id", "error", err)
		return
	}

	invitation := state.NewInvitation(invitationId)
	invitation.PassengerId = clientId
	invitation.DriverId = msg.DriverId
	invitation.CreatedAt = p.mainState.Clock.Now()

	p.mainState.GetInvitationsContainer().Add(invitation)

	p.driversHub.SendMessage(ws2.ClientMessage{
		ClientId: msg.DriverId,
		Message: messages.ServerOffersPassenger{
			PassengerId: clientId,
			Lat:         msg.Lat,
			Lon:         msg.Lon,
		},
	})
}

func (p *PassengersMessageHandler) handlePassengerCancelsInvite(clientId string, msg *messages.PassengerCancelsInvite) {
	p.logger.Info("passenger cancels invite", "client_id", clientId, "driver_id", msg.DriverId)

	invitations := p.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		PassengerId: clientId,
		DriverId:    msg.DriverId,
	})
	// the passenger can cancel a pending or accepted invitation; skip older finished ones with the same pair
	// TODO replace this loop with InvitationsFilter statuses once it accepts multiple values
	var invitation *state.Invitation
	for _, inv := range invitations {
		if inv.Status == state.InvitationStatusPending || inv.Status == state.InvitationStatusAccepted {
			invitation = inv
			break
		}
	}
	if invitation == nil {
		p.logger.Warn("PassengerCancelsInvite for unknown invitation", "client_id", clientId, "driver_id", msg.DriverId)
		return
	}

	p.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
		inv.Status = state.InvitationStatusCanceled
		inv.CancelledAt = p.mainState.Clock.Now()
	})

	p.driversHub.SendMessage(ws2.ClientMessage{
		ClientId: msg.DriverId,
		Message:  messages.ServerCancelsOffer{PassengerId: clientId},
	})
}

func (p *PassengersMessageHandler) handlePassengerRequestDriverCoords(clientId string, msg *messages.PassengerRequestDriverCoords) {
	p.logger.Info("passenger requests driver coordinates", "client_id", clientId, "driver_id", msg.DriverId)
	// TODO after MVP
}

// forcing to implement interface
var _ MessageHandler = &PassengersMessageHandler{}
