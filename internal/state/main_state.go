package state

import (
	"fmt"
	"sync"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/pkg/util"
)

type MainState struct {
	drivers     map[string]*Driver
	driversLock sync.Mutex

	passengers     map[string]*Passenger
	passengersLock sync.Mutex

	dedicatedEvents chan util.Event
}

func NewMainState() *MainState {
	return &MainState{
		drivers:         make(map[string]*Driver),
		passengers:      make(map[string]*Passenger),
		driversLock:     sync.Mutex{},
		passengersLock:  sync.Mutex{},
		dedicatedEvents: make(chan util.Event, 256),
	}
}

func (s *MainState) GetDriversSnapshot() []*Driver {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driversList := make([]*Driver, 0, len(s.drivers))
	for _, driver := range s.drivers {
		driversList = append(driversList, driver)
	}

	return driversList
}

func (s *MainState) AddDedicatedEvent(event util.Event) bool {
	select {
	case s.dedicatedEvents <- event:
		return true
	default:

		fmt.Println("dedicatedEvents channel is full or null")
		return false
	}
}

func (s *MainState) CreateRegistrationRelatedListener() func(event util.Event) {
	return func(event util.Event) {
		if !s.AddDedicatedEvent(event) {
			fmt.Println("dropped event, dedicatedEvents channel is full or null:", event.GetName())
		}
	}
}

func (s *MainState) HandleDedicatedEvents() {
	for {
		event := <-s.dedicatedEvents

		switch e := event.(type) {
		case *events.ClientRegisteredEvent:

			fmt.Println("TODO handle ClientRegisteredEvent: ", e.ClientId, e.ClientType)
		case *events.ClientUnregisteredEvent:

			fmt.Println("TODO handle ClientUnregisteredEvent: ", e.ClientId, e.ClientType)
		}

	}
}
