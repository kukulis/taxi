package handler

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
)

type DriverStatusChangeForDriverListener struct {
	MainState *state.MainState
}

func NewDriverStatusChangeForDriverListener(mainState *state.MainState) *DriverStatusChangeForDriverListener {
	return &DriverStatusChangeForDriverListener{
		MainState: mainState,
	}
}

func (l *DriverStatusChangeForDriverListener) Handle(e util.Event) {
	statusChangeEvent := e.(events.DriverStatusChangedEvent)
	//oldStatus := state.DriverStatus(statusChangeEvent.DriverStatusOld)
	newStatus := state.DriverStatus(statusChangeEvent.DriverStatusNew)

	l.MainState.UpdateDriver(statusChangeEvent.DriverId, func(driver *state.Driver) {
		driver.Status = newStatus
	})
}
