package handler

import (
	"log"
	"log/slog"
	"strconv"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

type PassengersMessageHandler struct {
	mainState                  *state.MainState
	driversHub                 ws2.HubInterface
	passengersHub              ws2.HubInterface
	logger                     *slog.Logger
	nominatimAddressWorkerPool *nominatim.AddressWorkerPool
}

func NewPassengersMessageHandler(
	mainState *state.MainState,
	driversHub ws2.HubInterface,
	passengersHub ws2.HubInterface,
	nominatimAddressWorkerPool *nominatim.AddressWorkerPool,
) *PassengersMessageHandler {
	return &PassengersMessageHandler{
		mainState:                  mainState,
		driversHub:                 driversHub,
		passengersHub:              passengersHub,
		logger:                     slog.Default().With("handler", "PassengersMessageHandler"),
		nominatimAddressWorkerPool: nominatimAddressWorkerPool,
	}
}

// HandleClientUnregistered marks a passenger offline when its client disconnects
func (p *PassengersMessageHandler) HandleClientUnregistered(clientId string, payload message_common.MessageInterface) {
	_ = payload.(*message_common.ClientUnregisteredMessage)

	p.mainState.UpdateClient(&state.Passenger{Id: clientId}, func(c state.Client) {
		passenger := c.(*state.Passenger)
		passenger.Online = false
		passenger.Active = false
	})
}

func (p *PassengersMessageHandler) HandleClientRegistered(clientId string, payload message_common.MessageInterface) {
	_ = payload.(*message_common.ClientRegisteredMessage)

	p.mainState.AddClient(&state.Passenger{
		Id:        clientId,
		CreatedAt: p.mainState.Clock.Now(),
	})

	// only the invitations still relevant to the passenger; rejected/cancelled/completed ones are done with
	invitations := p.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		PassengerId: clientId,
		Statuses: []string{
			state.InvitationStatusPending,
			state.InvitationStatusAccepted,
			state.InvitationStatusDriving,
		},
	})

	driverClients := p.mainState.SearchClientsByIds(util.ArrayMap(invitations, state.GetInvitationDriverId), state.ClientTypeDriver)
	drivers := util.ArrayMap(driverClients, func(c state.Client) *state.Driver { return c.(*state.Driver) })

	state.AssignDriversToInvitations(invitations, drivers)

	dtos := util.ArrayMap(invitations, func(invitation *state.Invitation) messages.InvitationDto {
		return messages.InvitationDto{
			InvitationId: invitation.Id,
			Status:       invitation.Status,
			DriverId:     invitation.DriverId,
			Lat:          strconv.FormatFloat(invitation.GetDriverLat(), 'f', -1, 64),
			Lon:          strconv.FormatFloat(invitation.GetDriverLon(), 'f', -1, 64),
		}
	})

	p.passengersHub.SendMessage(ws2.ClientMessage{
		ClientId: clientId,
		Message: messages.ServerRefreshPassengerInvitations{
			Invitations: dtos,
		},
	})
}

func (p *PassengersMessageHandler) HandleClientRespondsCoordinates(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientRespondsCoordinates)

	found := p.mainState.UpdateClient(&state.Passenger{Id: clientId}, func(c state.Client) {
		passenger := c.(*state.Passenger)
		passenger.Lat = msg.Lat
		passenger.Lon = msg.Lon
		passenger.CoordinatesReceivedAt = p.mainState.Clock.Now()
	})
	if !found {
		p.logger.Warn("ClientRespondsCoordinates from unknown passenger", "client_id", clientId)
		return
	}

	// ====================== send passenger's address =======================

	if p.nominatimAddressWorkerPool != nil {
		p.nominatimAddressWorkerPool.AddRequest(nominatim.NewAddressRequestToWorker(
			msg.Lat, msg.Lon, func(address string, err error) {
				if err != nil {
					log.Println("Failed to get address from nominatim", err.Error())
					return
				}

				truncatedAddress := util.TruncateAddress(address, 3)
				log.Println("received address: ", address, ". Truncated address: ", truncatedAddress)

				p.mainState.UpdateClient(&state.Passenger{Id: clientId}, func(c state.Client) {
					c.(*state.Passenger).Address = truncatedAddress
				})

				p.passengersHub.SendMessage(ws2.ClientMessage{
					ClientId: clientId,
					Message: messages.ServerSendsClientAddress{
						Address: truncatedAddress,
					},
				})
			},
		))
	}
}

func (p *PassengersMessageHandler) HandleClientRespondsInfo(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientRespondsInfo)

	found := p.mainState.UpdateClient(&state.Passenger{Id: clientId}, func(c state.Client) {
		passenger := c.(*state.Passenger)
		passenger.Phone = msg.Phone
		passenger.Name = msg.Name
		if msg.Ip != "" {
			passenger.Ip = msg.Ip
		}
		passenger.InfoReceivedAt = p.mainState.Clock.Now()
	})
	if !found {
		p.logger.Warn("ClientRespondsInfo from unknown passenger", "client_id", clientId)
	}
}

// HandleClientResponseError logs the passenger client's report that it didn't
// understand a request the server sent it. No MainState field maps to this,
// so it's just logged for now.
func (p *PassengersMessageHandler) HandleClientResponseError(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientResponseError)

	p.logger.Warn("client didn't understand a server request", "client_id", clientId, "error", msg.Error)
}

// HandleClientError logs the passenger client's own functionality error. No
// MainState field maps to this, so it's just logged for now.
func (p *PassengersMessageHandler) HandleClientError(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientError)

	p.logger.Error("client reported a functionality error", "client_id", clientId, "error", msg.Error)
}

func (p *PassengersMessageHandler) HandlePassengerInvitesDriver(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.PassengerInvitesDriver)

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

func (p *PassengersMessageHandler) HandlePassengerCancelsInvite(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.PassengerCancelsInvite)

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

func (p *PassengersMessageHandler) HandlePassengerRequestDriverCoords(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.PassengerRequestDriverCoords)

	p.logger.Info("passenger requests driver coordinates", "client_id", clientId, "driver_id", msg.DriverId)
	// TODO after MVP
}
