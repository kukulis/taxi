package handler

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
)

func TestDriverStatusChangeForDriverListener_UpdatesDriverStatus(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	dispatcher := util.NewDispatcher()
	listener := NewDriverStatusChangeForDriverListener(mainState)
	dispatcher.AddListener(events.DriverStatusChangedEventName, listener.Handle)

	dispatcher.Dispatch(events.DriverStatusChangedEvent{
		DriverId:        driverId,
		DriverStatusOld: string(state.DriverStatusIdle),
		DriverStatusNew: string(state.DriverStatusWorking),
	})

	driver := findDriver(mainState, driverId)
	if driver == nil {
		t.Fatalf("findDriver(%q) = nil, want a driver", driverId)
	}
	if driver.Status != state.DriverStatusWorking {
		t.Errorf("driver.Status = %q, want %q", driver.Status, state.DriverStatusWorking)
	}
}

func TestDriverStatusChangeForDriverListener_UnknownDriverDoesNotPanic(t *testing.T) {
	mainState := state.NewMainState(util.NewFixedClock(time.Now()))

	listener := NewDriverStatusChangeForDriverListener(mainState)

	listener.Handle(events.DriverStatusChangedEvent{
		DriverId:        "unknown-driver",
		DriverStatusOld: string(state.DriverStatusIdle),
		DriverStatusNew: string(state.DriverStatusWorking),
	})

	if driver := findDriver(mainState, "unknown-driver"); driver != nil {
		t.Errorf("findDriver(%q) = %+v, want nil (no driver should have been created)", "unknown-driver", driver)
	}
}
