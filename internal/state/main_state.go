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
		drivers:        make(map[string]*Driver),
		passengers:     make(map[string]*Passenger),
		driversLock:    sync.Mutex{},
		passengersLock: sync.Mutex{},
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

func (s *MainState) AddDedicatedEvent(event util.Event) {
	s.dedicatedEvents <- event
}

func (s *MainState) CreateRegistrationRelatedListener() func(event util.Event) {
	return func(event util.Event) {
		s.dedicatedEvents <- event
	}
}

func (s *MainState) HandleDedicatedEvents() {
	for {
		event := <-s.dedicatedEvents

		switch e := event.(type) {
		case *events.DriverRegisteredEvent:

			// TODO
			fmt.Println(e.ClientId)
		case *events.DriverUnregisteredEvent:

			// TODO
			fmt.Println(e.ClientId)
		case *events.PassengerRegisteredEvent:

			//TODO
			fmt.Println(e.ClientId)
		case *events.PassengerUnregisteredEvent:

			// TODO
			fmt.Println(e.ClientId)
		}

	}
}
