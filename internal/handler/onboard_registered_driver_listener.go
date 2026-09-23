package handler

import (
	"fmt"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
)

// OnboardRegisteredDriverListener onboards a newly registered driver.
// Dedicated to take a registration channel from MainState
type OnboardRegisteredDriverListener struct {
	mainState *state.MainState
}

func NewOnboardRegisteredDriverListener(mainState *state.MainState) *OnboardRegisteredDriverListener {
	return &OnboardRegisteredDriverListener{
		mainState: mainState,
	}
}

func (l *OnboardRegisteredDriverListener) Handle(e util.Event) {
	registeredEvent := e.(*events.ClientRegisteredEvent)

	if registeredEvent.ClientType != events.ClientTypeDriver {
		return
	}

	// TODO move to channel later: this calls CreateDriver directly, synchronously,
	// in the Hub.Run() goroutine (Dispatcher.Dispatch calls listeners synchronously
	// in the dispatching goroutine). CreateDriver takes MainState's driversLock
	// mutex, so lock contention here can slow down the Hub's own loop - exactly
	// what the dedicatedEvents channel exists to avoid. Stopgap until this
	// listener gets its own channel/consumer, per the earlier MainState refactor plan.
	l.mainState.CreateDriver(registeredEvent.ClientId)

	// TODO after MVP we request driver coordinates, but we do that using double-step event after registering a
	// a driver in a MainState in the registrations channel consumer
	fmt.Println("TODO onboard driver")
}
