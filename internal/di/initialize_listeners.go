package di

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/pkg/util"
)

func InitializeListeners(dispatcher *util.Dispatcher) {
	dispatcher.AddListener(events.DriverRegisteredEventName, GetMainState().CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.DriverUnregisteredEventName, GetMainState().CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.PassengerRegisteredEventName, GetMainState().CreateRegistrationRelatedListener())
	dispatcher.AddListener(events.PassengerUnregisteredEventName, GetMainState().CreateRegistrationRelatedListener())
}
