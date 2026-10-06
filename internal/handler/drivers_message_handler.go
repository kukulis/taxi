package handler

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
)

type DriversMessageHandler struct {
	mainState                  *state.MainState
	driversHub                 ws.HubInterface
	passengersHub              ws.HubInterface
	logger                     *slog.Logger
	dispatcher                 *util.Dispatcher
	nominatimAddressWorkerPool *nominatim.AddressWorkerPool
}

func NewDriversMessageHandler(
	mainState *state.MainState,
	driversHub ws.HubInterface,
	passengersHub ws.HubInterface,
	dispatcher *util.Dispatcher,
	nominatimAddressWorkerPool *nominatim.AddressWorkerPool,
) *DriversMessageHandler {
	return &DriversMessageHandler{
		mainState:                  mainState,
		driversHub:                 driversHub,
		passengersHub:              passengersHub,
		logger:                     slog.Default().With("handler", "DriversMessageHandler"),
		dispatcher:                 dispatcher,
		nominatimAddressWorkerPool: nominatimAddressWorkerPool,
	}
}

func (d *DriversMessageHandler) Handle(messageChannel <-chan ws.ClientMessage) {

	for message := range messageChannel {

		switch msg := message.Message.(type) {
		case *messages.ClientRespondsCoordinates:
			d.HandleClientRespondsCoordinates(message.ClientId, msg)
		case *messages.ClientRespondsInfo:
			d.HandleClientRespondsInfo(message.ClientId, msg)
		case *messages.ClientResponseError:
			d.HandleClientResponseError(message.ClientId, msg)
		case *messages.ClientError:
			d.HandleClientError(message.ClientId, msg)
		case *messages.DriverCancelsOffer:
			d.HandleDriverCancelsOffer(message.ClientId, msg)
		case *messages.DriverAcceptsOffer:
			d.HandleDriverAcceptsOffer(message.ClientId, msg)
		case *messages.DriverRejectsOffer:
			d.HandleDriverRejectsOffer(message.ClientId, msg)
		case *messages.DriverChangesStatus:
			d.HandleDriverChangesStatus(message.ClientId, msg)
		case *message_common.ClientRegisteredMessage:
			d.HandleClientRegistered(message.ClientId, msg)
		case *message_common.ClientUnregisteredMessage:
			d.HandleClientUnregistered(message.ClientId, msg)
		default:
			d.logger.Warn("unexpected message type",
				"client_id", message.ClientId,
				"message_type", fmt.Sprintf("%T", message.Message),
			)
		}
	}
}

// HandleClientUnregistered marks a driver offline when its client disconnects
func (d *DriversMessageHandler) HandleClientUnregistered(clientId string, payload message_common.MessageInterface) {
	_ = payload.(*message_common.ClientUnregisteredMessage)

	d.mainState.UpdateClient(&state.Driver{Id: clientId}, func(c state.Client) {
		driver := c.(*state.Driver)
		driver.Online = false
		driver.Active = false
	})
}

// HandleClientRegistered handling virtual message about a driver registration
func (d *DriversMessageHandler) HandleClientRegistered(clientId string, payload message_common.MessageInterface) {
	_ = payload.(*message_common.ClientRegisteredMessage)

	d.mainState.AddClient(&state.Driver{
		Id:        clientId,
		CreatedAt: d.mainState.Clock.Now(),
		Status:    state.DriverStatusResting,
	})

	driverFilter := state.NewEmptyFilter()
	driverFilter.Id = clientId
	driverFilter.Type = string(state.ClientTypeDriver)
	driver := d.mainState.SearchClients(driverFilter)[0].(*state.Driver)

	// ============  Refresh driver status ===================================

	d.driversHub.SendMessage(ws.ClientMessage{
		ClientId: clientId,
		Message: messages.ServerRefreshDriverStatus{
			Status: string(driver.Status),
		},
	})

	// ============ Send invitations related to the driver ====================

	// only the offers still relevant to the driver; rejected/cancelled/completed ones are done with
	invitations := d.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		DriverId: clientId,
		Statuses: []string{
			state.InvitationStatusPending,
			state.InvitationStatusAccepted,
			state.InvitationStatusDriving,
		},
	})

	passengerClients := d.mainState.SearchClientsByIds(util.ArrayMap(invitations, state.GetInvitationPassengerId), state.ClientTypePassenger)
	passengers := util.ArrayMap(passengerClients, func(c state.Client) *state.Passenger { return c.(*state.Passenger) })

	state.AssignPassengersToInvitations(invitations, passengers)

	offers := util.ArrayMap(invitations, func(invitation *state.Invitation) messages.OfferDto {
		return messages.OfferDto{
			InvitationId: invitation.Id,
			Status:       invitation.Status,
			PassengerId:  invitation.PassengerId,
			Lat:          strconv.FormatFloat(invitation.GetPassengerLat(), 'f', -1, 64),
			Lon:          strconv.FormatFloat(invitation.GetPassengerLon(), 'f', -1, 64),
		}
	})

	d.driversHub.SendMessage(ws.ClientMessage{
		ClientId: clientId,
		Message: messages.ServerRefreshDriverOffers{
			Offers: offers,
		},
	})

	// ====================== finito ======================================
}

// HandleClientRespondsCoordinates write received coordinates to a main state
func (d *DriversMessageHandler) HandleClientRespondsCoordinates(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientRespondsCoordinates)

	found := d.mainState.UpdateClient(&state.Driver{Id: clientId}, func(c state.Client) {
		driver := c.(*state.Driver)
		driver.Lat = msg.Lat
		driver.Lon = msg.Lon
		driver.CoordinatesReceivedAt = d.mainState.Clock.Now()
	})
	if !found {
		d.logger.Warn("ClientRespondsCoordinates from unknown driver", "client_id", clientId)
		return
	}

	// ====================== send driver's address =======================

	if d.nominatimAddressWorkerPool != nil {
		d.nominatimAddressWorkerPool.AddRequest(nominatim.NewAddressRequestToWorker(
			msg.Lat, msg.Lon, func(address string, err error) {
				if err != nil {
					log.Println("Failed to get address from nominatim", err.Error())
					return
				}
				truncatedAddress := util.TruncateAddress(address, 3)
				log.Println("received address: ", address, ". Truncated address: ", truncatedAddress)

				d.mainState.UpdateClient(&state.Driver{Id: clientId}, func(c state.Client) {
					c.(*state.Driver).Address = truncatedAddress
				})

				d.driversHub.SendMessage(ws.ClientMessage{
					ClientId: clientId,
					Message: messages.ServerSendsClientAddress{
						Address: truncatedAddress,
					},
				})
			},
		))
	}
}

func (d *DriversMessageHandler) HandleClientRespondsInfo(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientRespondsInfo)

	found := d.mainState.UpdateClient(&state.Driver{Id: clientId}, func(c state.Client) {
		driver := c.(*state.Driver)
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

// HandleClientResponseError logs the driver client's report that it didn't
// understand a request the server sent it. No MainState field maps to this,
// so it's just logged for now.
func (d *DriversMessageHandler) HandleClientResponseError(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientResponseError)

	d.logger.Warn("client didn't understand a server request", "client_id", clientId, "error", msg.Error)
}

// HandleClientError logs the driver client's own functionality error (e.g. it
// can't read device coordinates). No MainState field maps to this, so it's
// just logged for now.
func (d *DriversMessageHandler) HandleClientError(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.ClientError)

	d.logger.Error("client reported a functionality error", "client_id", clientId, "error", msg.Error)
	// TODO set client Active=false
}

func (d *DriversMessageHandler) HandleDriverCancelsOffer(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.DriverCancelsOffer)

	d.logger.Info("driver cancels offer", "client_id", clientId, "passenger_id", msg.PassengerId)

	// the driver's "cancel" button is only shown on an accepted offer
	invitations := d.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		DriverId:    clientId,
		PassengerId: msg.PassengerId,
		Status:      state.InvitationStatusAccepted,
	})
	if len(invitations) == 0 {
		d.logger.Warn("DriverCancelsOffer for unknown invitation", "client_id", clientId, "passenger_id", msg.PassengerId)
		return
	}
	invitation := invitations[0]

	d.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
		inv.Status = state.InvitationStatusCanceled
		inv.CancelledAt = d.mainState.Clock.Now()
	})

	d.passengersHub.SendMessage(ws.ClientMessage{
		ClientId: invitation.PassengerId,
		Message:  messages.ServerNotifiesInviteCanceled{DriverId: clientId},
	})
}

func (d *DriversMessageHandler) HandleDriverAcceptsOffer(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.DriverAcceptsOffer)

	d.logger.Info("driver accepts offer", "client_id", clientId, "passenger_id", msg.PassengerId)

	// only a pending invitation can be accepted; older rejected/completed ones with the same pair must be skipped
	invitations := d.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		DriverId:    clientId,
		PassengerId: msg.PassengerId,
		Status:      state.InvitationStatusPending,
	})
	if len(invitations) == 0 {
		d.logger.Warn("DriverAcceptsOffer for unknown invitation", "client_id", clientId, "passenger_id", msg.PassengerId)
		return
	}
	accepted := invitations[0]

	d.mainState.GetInvitationsContainer().Update(accepted.Id, func(inv *state.Invitation) {
		inv.Status = state.InvitationStatusAccepted
		inv.AcceptedAt = d.mainState.Clock.Now()
	})

	// accepting one invite drops every other pending invite for this driver
	others := d.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		DriverId: clientId,
		Status:   state.InvitationStatusPending,
	})
	for _, other := range others {
		d.mainState.GetInvitationsContainer().Update(other.Id, func(inv *state.Invitation) {
			inv.Status = state.InvitationStatusRejected
			inv.RejectedAt = d.mainState.Clock.Now()
		})
		d.passengersHub.SendMessage(ws.ClientMessage{
			ClientId: other.PassengerId,
			Message:  messages.ServerNotifiesInviteRejected{DriverId: clientId},
		})
	}

	d.passengersHub.SendMessage(ws.ClientMessage{
		ClientId: accepted.PassengerId,
		Message:  messages.ServerNotifiesInviteAccepted{DriverId: clientId},
	})
}

func (d *DriversMessageHandler) HandleDriverRejectsOffer(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.DriverRejectsOffer)

	d.logger.Info("driver rejects offer", "client_id", clientId, "passenger_id", msg.PassengerId)

	invitations := d.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{
		DriverId:    clientId,
		PassengerId: msg.PassengerId,
		Status:      state.InvitationStatusPending,
	})
	if len(invitations) == 0 {
		d.logger.Warn("DriverRejectsOffer for unknown invitation", "client_id", clientId, "passenger_id", msg.PassengerId)
		return
	}
	invitation := invitations[0]

	d.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
		inv.Status = state.InvitationStatusRejected
		inv.RejectedAt = d.mainState.Clock.Now()
	})

	d.passengersHub.SendMessage(ws.ClientMessage{
		ClientId: invitation.PassengerId,
		Message:  messages.ServerNotifiesInviteRejected{DriverId: clientId},
	})
}

func (d *DriversMessageHandler) HandleDriverChangesStatus(clientId string, payload message_common.MessageInterface) {
	msg := payload.(*messages.DriverChangesStatus)

	d.logger.Info("driver changes status", "client_id", clientId, "status", msg.Status)

	filter := state.NewEmptyFilter()
	filter.Id = clientId
	filter.Type = string(state.ClientTypeDriver)
	found := d.mainState.SearchClients(filter)

	if len(found) == 0 {
		d.logger.Warn("DriverChangesStatus for unknown driver", "client_id", clientId)
		return
	}
	driver := found[0].(*state.Driver)

	event := events.DriverStatusChangedEvent{
		DriverId:        clientId,
		DriverStatusOld: string(driver.Status),
		DriverStatusNew: string(msg.Status),
	}

	// maybe the switch is a bit wasteful now but ok
	switch msg.Status {
	case state.DriverStatusIdle:
		d.dispatcher.Dispatch(event)
	case state.DriverStatusWorking:
		d.dispatcher.Dispatch(event)
	case state.DriverStatusOffline:
		d.dispatcher.Dispatch(event)
	case state.DriverStatusResting:
		d.dispatcher.Dispatch(event)
	default:
		d.logger.Warn("driver changed to an unknown status", "client_id", clientId, "status", msg.Status)
	}
}

// forcing to implement interface
var _ MessageHandler = &DriversMessageHandler{}
