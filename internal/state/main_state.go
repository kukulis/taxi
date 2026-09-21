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
		driverCopy := *driver
		driversList = append(driversList, &driverCopy)
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
			switch e.ClientType {
			case events.ClientTypeDriver:
				s.handleDriverRegistered(e)
			case events.ClientTypePassenger:
				s.handlePassengerRegistered(e)
			}
		case *events.ClientUnregisteredEvent:
			switch e.ClientType {
			case events.ClientTypeDriver:
				s.handleDriverUnregistered(e)
			case events.ClientTypePassenger:
				s.handlePassengerUnregistered(e)
			}
		}

	}
}

func (s *MainState) handleDriverRegistered(e *events.ClientRegisteredEvent) {

	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driver, ok := s.drivers[e.ClientId]

	if !ok {
		driver = NewDriver()
		s.drivers[e.ClientId] = driver
		driver.Id = e.ClientId
	}

	driver.Active = true

	fmt.Println("Driver registered:", e.ClientId)
}

func (s *MainState) handleDriverUnregistered(e *events.ClientUnregisteredEvent) {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driver, ok := s.drivers[e.ClientId]

	if ok {
		driver.Active = false
	}
	fmt.Println("Driver unregistered/inactivated:", e.ClientId)
}

func (s *MainState) handlePassengerRegistered(e *events.ClientRegisteredEvent) {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()
	fmt.Println("TODO handle passenger registered:", e.ClientId)
}

func (s *MainState) handlePassengerUnregistered(e *events.ClientUnregisteredEvent) {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()
	fmt.Println("TODO handle passenger unregistered:", e.ClientId)
}
