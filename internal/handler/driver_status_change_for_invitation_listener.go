package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
)

// DriverStatusChangeForInvitationListener The listener is moved outside 'state' package to avoid cyclic dependency
type DriverStatusChangeForInvitationListener struct {
	mainState     *state.MainState
	passengersHub ws.HubInterface
}

func NewDriverStatusChangeForInvitationListener(
	mainState *state.MainState,
	passengersHub ws.HubInterface,
) *DriverStatusChangeForInvitationListener {
	return &DriverStatusChangeForInvitationListener{
		mainState:     mainState,
		passengersHub: passengersHub,
	}
}

func (l *DriverStatusChangeForInvitationListener) Handle(e util.Event) {
	fmt.Println("Handle status change for invitation TODO")

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
			// TODO send message to passenger
			l.passengersHub.SendMessage(
				ws.ClientMessage{
					ClientId: invitation.PassengerId,
					Message: messages.ServerNotifiesVoyageFinished{
						DriverId:    statusChangeEvent.DriverId,
						PassengerId: invitation.PassengerId,
					},
				})
		}

		return
	}

	if oldStatus == state.DriverStatusIdle && newStatus == state.DriverStatusWorking {
		invitations := l.mainState.GetInvitationsContainer().GetInvitationsSnapshot(
			&state.InvitationsFilter{
				DriverId: statusChangeEvent.DriverId,
				Status:   state.InvitationStatusAccepted,
			})
		if len(invitations) > 0 {
			invitation := invitations[0]
			// TODO decide what to do if there are more invitations
			l.mainState.GetInvitationsContainer().Update(invitation.Id, func(inv *state.Invitation) {
				inv.Status = state.InvitationStatusDriving
			})

			// TODO send message to passenger
			l.passengersHub.SendMessage(
				ws.ClientMessage{
					ClientId: invitation.PassengerId,
					Message: messages.ServerNotifiesVoyageStarted{
						DriverId:    statusChangeEvent.DriverId,
						PassengerId: invitation.PassengerId,
					},
				})
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
