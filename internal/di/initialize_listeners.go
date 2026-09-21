package di

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
)

func InitializeListenersFromMainState(mainState *state.MainState, dispatcher *util.Dispatcher) {
	dispatcher.AddListener(events.ClientRegisteredEventName, mainState.CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.ClientUnregisteredEventName, mainState.CreateRegistrationRelatedListener())
}
