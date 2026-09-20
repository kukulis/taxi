package di

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
)

func InitializeListenersFromMainState(mainState *state.MainState, dispatcher *util.Dispatcher) {
	dispatcher.AddListener(events.DriverRegisteredEventName, mainState.CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.DriverUnregisteredEventName, mainState.CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.PassengerRegisteredEventName, mainState.CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.PassengerUnregisteredEventName, mainState.CreateRegistrationRelatedListener())
}
