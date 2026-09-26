package handler

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

func TestDriverStatusChangeForInvitationListener_StartsVoyageOnIdleToWorking(t *testing.T) {
	const driverId = "driver-1"
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.CreateDriver(driverId)
	mainState.CreatePassenger(passengerId)

	invitation := state.NewInvitation("inv-1")
	invitation.DriverId = driverId
	invitation.PassengerId = passengerId
	invitation.Status = state.InvitationStatusAccepted
	mainState.GetInvitationsContainer().Add(invitation)

	var sent []ws2.ClientMessage
	passengersHub := ws2.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws2.ClientMessage) {
		sent = append(sent, msg)
	}

	dispatcher := util.NewDispatcher()
	listener := NewDriverStatusChangeForInvitationListener(mainState, ws2.NewHubMock(), passengersHub)
	dispatcher.AddListener(events.DriverStatusChangedEventName, listener.Handle)

	dispatcher.Dispatch(events.DriverStatusChangedEvent{
		DriverId:        driverId,
		DriverStatusOld: string(state.DriverStatusIdle),
		DriverStatusNew: string(state.DriverStatusWorking),
	})

	updated := mainState.GetInvitationsContainer().GetSnapshot(invitation.Id)
	if updated == nil {
		t.Fatalf("invitation %q disappeared from the container", invitation.Id)
	}
	if updated.Status != state.InvitationStatusDriving {
		t.Errorf("invitation.Status = %q, want %q", updated.Status, state.InvitationStatusDriving)
	}

	if len(sent) != 1 {
		t.Fatalf("passengersHub received %d messages, want 1: %+v", len(sent), sent)
	}
	if sent[0].ClientId != passengerId {
		t.Errorf("message ClientId = %q, want %q", sent[0].ClientId, passengerId)
	}
	started, ok := sent[0].Message.(messages.ServerNotifiesVoyageStarted)
	if !ok {
		t.Fatalf("message type = %T, want messages.ServerNotifiesVoyageStarted", sent[0].Message)
	}
	if started.DriverId != driverId || started.PassengerId != passengerId {
		t.Errorf("ServerNotifiesVoyageStarted = %+v, want DriverId=%q PassengerId=%q", started, driverId, passengerId)
	}
}

func TestDriverStatusChangeForInvitationListener_FinishesVoyageOnWorkingToIdle(t *testing.T) {
	const driverId = "driver-1"
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.CreateDriver(driverId)
	mainState.CreatePassenger(passengerId)

	invitation := state.NewInvitation("inv-1")
	invitation.DriverId = driverId
	invitation.PassengerId = passengerId
	invitation.Status = state.InvitationStatusDriving
	mainState.GetInvitationsContainer().Add(invitation)

	var sent []ws2.ClientMessage
	passengersHub := ws2.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws2.ClientMessage) {
		sent = append(sent, msg)
	}

	dispatcher := util.NewDispatcher()
	listener := NewDriverStatusChangeForInvitationListener(mainState, ws2.NewHubMock(), passengersHub)
	dispatcher.AddListener(events.DriverStatusChangedEventName, listener.Handle)

	dispatcher.Dispatch(events.DriverStatusChangedEvent{
		DriverId:        driverId,
		DriverStatusOld: string(state.DriverStatusWorking),
		DriverStatusNew: string(state.DriverStatusIdle),
	})

	updated := mainState.GetInvitationsContainer().GetSnapshot(invitation.Id)
	if updated == nil {
		t.Fatalf("invitation %q disappeared from the container", invitation.Id)
	}
	if updated.Status != state.InvitationStatusCompleted {
		t.Errorf("invitation.Status = %q, want %q", updated.Status, state.InvitationStatusCompleted)
	}

	if len(sent) != 1 {
		t.Fatalf("passengersHub received %d messages, want 1: %+v", len(sent), sent)
	}
	if sent[0].ClientId != passengerId {
		t.Errorf("message ClientId = %q, want %q", sent[0].ClientId, passengerId)
	}
	finished, ok := sent[0].Message.(messages.ServerNotifiesVoyageFinished)
	if !ok {
		t.Fatalf("message type = %T, want messages.ServerNotifiesVoyageFinished", sent[0].Message)
	}
	if finished.DriverId != driverId || finished.PassengerId != passengerId {
		t.Errorf("ServerNotifiesVoyageFinished = %+v, want DriverId=%q PassengerId=%q", finished, driverId, passengerId)
	}
}

func TestDriverStatusChangeForInvitationListener_CancelsPendingInvitationsOnResting(t *testing.T) {
	const driverId = "driver-1"
	const otherDriverId = "driver-2"
	const passenger1Id = "passenger-1"
	const passenger2Id = "passenger-2"
	const otherPassengerId = "passenger-3"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.CreateDriver(driverId)
	mainState.CreateDriver(otherDriverId)
	mainState.CreatePassenger(passenger1Id)
	mainState.CreatePassenger(passenger2Id)
	mainState.CreatePassenger(otherPassengerId)

	invitation1 := state.NewInvitation("inv-1")
	invitation1.DriverId = driverId
	invitation1.PassengerId = passenger1Id
	invitation1.Status = state.InvitationStatusPending
	mainState.GetInvitationsContainer().Add(invitation1)

	invitation2 := state.NewInvitation("inv-2")
	invitation2.DriverId = driverId
	invitation2.PassengerId = passenger2Id
	invitation2.Status = state.InvitationStatusPending
	mainState.GetInvitationsContainer().Add(invitation2)

	// unrelated invitation, different driver - must not be touched
	otherInvitation := state.NewInvitation("inv-3")
	otherInvitation.DriverId = otherDriverId
	otherInvitation.PassengerId = otherPassengerId
	otherInvitation.Status = state.InvitationStatusPending
	mainState.GetInvitationsContainer().Add(otherInvitation)

	var sent []ws2.ClientMessage
	passengersHub := ws2.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws2.ClientMessage) {
		sent = append(sent, msg)
	}

	dispatcher := util.NewDispatcher()
	listener := NewDriverStatusChangeForInvitationListener(mainState, ws2.NewHubMock(), passengersHub)
	dispatcher.AddListener(events.DriverStatusChangedEventName, listener.Handle)

	dispatcher.Dispatch(events.DriverStatusChangedEvent{
		DriverId:        driverId,
		DriverStatusOld: string(state.DriverStatusIdle),
		DriverStatusNew: string(state.DriverStatusResting),
	})

	for _, id := range []string{invitation1.Id, invitation2.Id} {
		updated := mainState.GetInvitationsContainer().GetSnapshot(id)
		if updated == nil {
			t.Fatalf("invitation %q disappeared from the container", id)
		}
		if updated.Status != state.InvitationStatusCanceled {
			t.Errorf("invitation %q Status = %q, want %q", id, updated.Status, state.InvitationStatusCanceled)
		}
	}

	untouched := mainState.GetInvitationsContainer().GetSnapshot(otherInvitation.Id)
	if untouched == nil {
		t.Fatalf("invitation %q disappeared from the container", otherInvitation.Id)
	}
	if untouched.Status != state.InvitationStatusPending {
		t.Errorf("unrelated invitation %q Status = %q, want unchanged %q", otherInvitation.Id, untouched.Status, state.InvitationStatusPending)
	}

	if len(sent) != 2 {
		t.Fatalf("passengersHub received %d messages, want 2: %+v", len(sent), sent)
	}

	notifiedPassengers := map[string]bool{}
	for _, msg := range sent {
		notifiedPassengers[msg.ClientId] = true
		canceled, ok := msg.Message.(messages.ServerNotifiesInviteCanceled)
		if !ok {
			t.Fatalf("message type = %T, want messages.ServerNotifiesInviteCanceled", msg.Message)
		}
		if canceled.DriverId != driverId {
			t.Errorf("ServerNotifiesInviteCanceled.DriverId = %q, want %q", canceled.DriverId, driverId)
		}
	}
	for _, wantPassenger := range []string{passenger1Id, passenger2Id} {
		if !notifiedPassengers[wantPassenger] {
			t.Errorf("passenger %q was never notified, notified = %v", wantPassenger, notifiedPassengers)
		}
	}
}
