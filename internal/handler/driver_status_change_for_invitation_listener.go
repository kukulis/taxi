package handler

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
	"github.com/bytedance/gopkg/util/logger"
)

// DriverStatusChangeForInvitationListener The listener is moved outside 'state' package to avoid cyclic dependency
type DriverStatusChangeForInvitationListener struct {
	mainState     *state.MainState
	driversHub    ws.HubInterface
	passengersHub ws.HubInterface
}

func NewDriverStatusChangeForInvitationListener(
	mainState *state.MainState,
	driversHub ws.HubInterface,
	passengersHub ws.HubInterface,
) *DriverStatusChangeForInvitationListener {
	return &DriverStatusChangeForInvitationListener{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
	}
}

func (l *DriverStatusChangeForInvitationListener) Handle(e util.Event) {
	logger.Info("DriverStatusChangeForInvitationListener.Handle called with event ", e.GetName())

	statusChangeEvent := e.(events.DriverStatusChangedEvent)
	oldStatus := state.DriverStatus(statusChangeEvent.DriverStatusOld)
	newStatus := state.DriverStatus(statusChangeEvent.DriverStatusNew)

	// we will decide offline status change in a separate block
	// The offline status does not mean driver is stopped, so nothing to change
	if oldStatus == state.DriverStatusWorking && newStatus != state.DriverStatusOffline {
		invitations := l.mainState.GetInvitationsContainer().GetInvitationsSnapshot(
			&state.InvitationsFilter{
				DriverId: statusChangeEvent.DriverId,
				Status:   state.InvitationStatusDriving,
			})
		if len(invitations) > 0 {
			invitation := invitations[0]

			// TODO if there are other invitations
			l.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
				inv.Status = state.InvitationStatusCompleted
			})
			voyageFinished := messages.ServerNotifiesVoyageFinished{
				DriverId:    statusChangeEvent.DriverId,
				PassengerId: invitation.PassengerId,
			}
			l.passengersHub.SendMessage(ws.ClientMessage{ClientId: invitation.PassengerId, Message: voyageFinished})
			l.driversHub.SendMessage(ws.ClientMessage{ClientId: statusChangeEvent.DriverId, Message: voyageFinished})
		}

		return
	}

	if oldStatus == state.DriverStatusIdle && newStatus == state.DriverStatusWorking {
		invitations := l.mainState.GetInvitationsContainer().GetInvitationsSnapshot(
			&state.InvitationsFilter{
				DriverId: statusChangeEvent.DriverId,
				Status:   state.InvitationStatusAccepted,
			})
		logger.Info("DriverStatusChangeForInvitationListener.Handle: idle->working, invitations count: ", len(invitations), "")
		if len(invitations) > 0 {
			invitation := invitations[0]
			// TODO decide what to do if there are more invitations
			l.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
				inv.Status = state.InvitationStatusDriving
			})

			voyageStarted := messages.ServerNotifiesVoyageStarted{
				DriverId:    statusChangeEvent.DriverId,
				PassengerId: invitation.PassengerId,
			}
			l.passengersHub.SendMessage(ws.ClientMessage{ClientId: invitation.PassengerId, Message: voyageStarted})
			l.driversHub.SendMessage(ws.ClientMessage{ClientId: statusChangeEvent.DriverId, Message: voyageStarted})
		}

		return
	}

	if newStatus == state.DriverStatusResting {
		pending := l.mainState.GetInvitationsContainer().GetInvitationsSnapshot(
			&state.InvitationsFilter{
				DriverId: statusChangeEvent.DriverId,
				Status:   state.InvitationStatusPending,
			})
		for _, invitation := range pending {
			l.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
				inv.Status = state.InvitationStatusCanceled
				inv.CancelledAt = l.mainState.Clock.Now()
			})

			l.passengersHub.SendMessage(
				ws.ClientMessage{
					ClientId: invitation.PassengerId,
					Message:  messages.ServerNotifiesInviteCanceled{DriverId: statusChangeEvent.DriverId},
				})
		}

		return
	}

	// TODO other statuses changes related to invitations will be solved after planning ( probably after MVP )
}
