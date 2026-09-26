package state

import (
	"fmt"
	"sort"
	"sync"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/pkg/util"
	"github.com/bytedance/gopkg/util/logger"
)

type MainState struct {
	drivers     map[string]*Driver
	driversLock sync.Mutex

	passengers     map[string]*Passenger
	passengersLock sync.Mutex

	// TODO move to the OnboardRegisteredDriverListener
	dedicatedEvents chan util.Event

	// Clock provides the current time for timestamps such as Driver/Passenger
	// CreatedAt.
	Clock util.Clock

	invitationsContainer *InvitationsContainer
}

func NewMainState(clock util.Clock) *MainState {
	return &MainState{
		drivers:              make(map[string]*Driver),
		passengers:           make(map[string]*Passenger),
		driversLock:          sync.Mutex{},
		passengersLock:       sync.Mutex{},
		dedicatedEvents:      make(chan util.Event, 256),
		Clock:                clock,
		invitationsContainer: NewInvitationsContainer(),
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

func (s *MainState) GetDriverById(id string) *Driver {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driver, ok := s.drivers[id]
	if !ok {
		return nil
	}

	driverCopy := *driver
	return &driverCopy
}

// UpdateDriver locks the drivers map, looks up the driver by id and, if found,
// calls mutate on it in place. It reports whether the driver was found.
func (s *MainState) UpdateDriver(id string, mutate func(*Driver)) bool {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driver, ok := s.drivers[id]
	if !ok {
		return false
	}

	mutate(driver)
	return true
}

// CreateDriver adds a new driver with the given id, or reactivates the existing
// record if one is already present, and returns it.
func (s *MainState) CreateDriver(id string) *Driver {

	logger.Info("MainState.CreateDriver called")
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driver, ok := s.drivers[id]
	if !ok {
		driver = NewDriver()
		driver.Id = id
		driver.CreatedAt = s.Clock.Now()
		s.drivers[id] = driver
		driver.Status = DriverStatusResting
	}

	driver.Active = true
	return driver
}

// DriverDistance pairs a Driver with its distance in kilometers from a query point.
type DriverDistance struct {
	Driver   *Driver
	Distance float64
}

// GetNearestDrivers returns every active driver with known coordinates
// (CoordinatesReceivedAt set), sorted by ascending distance in kilometers from
// (lat, lon). Brute force: computes the distance to every active driver, no
// spatial index. Callers decide how many results to actually use.
func (s *MainState) GetNearestDrivers(lat, lon float64) []DriverDistance {
	drivers := s.GetDriversSnapshot()

	result := make([]DriverDistance, 0, len(drivers))
	for _, driver := range drivers {
		if !driver.Active || driver.CoordinatesReceivedAt.IsZero() {
			continue
		}
		result = append(result, DriverDistance{
			Driver:   driver,
			Distance: util.HaversineDistanceKm(lat, lon, driver.Lat, driver.Lon),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})

	return result
}

// RemoveDriver marks the driver with the given id inactive. The record is kept,
// not deleted, so it stays visible in GetDriversSnapshot/GetDriverById.
func (s *MainState) RemoveDriver(id string) {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	if driver, ok := s.drivers[id]; ok {
		driver.Active = false
	}
}

func (s *MainState) GetPassengersSnapshot() []*Passenger {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	passengersList := make([]*Passenger, 0, len(s.passengers))
	for _, passenger := range s.passengers {
		passengerCopy := *passenger
		passengersList = append(passengersList, &passengerCopy)
	}

	return passengersList
}

// UpdatePassenger locks the passengers map, looks up the passenger by id and, if found,
// calls mutate on it in place. It reports whether the passenger was found.
func (s *MainState) UpdatePassenger(id string, mutate func(*Passenger)) bool {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	passenger, ok := s.passengers[id]
	if !ok {
		return false
	}

	mutate(passenger)
	return true
}

// CreatePassenger adds a new passenger with the given id, or reactivates the existing
// record if one is already present, and returns it.
func (s *MainState) CreatePassenger(id string) *Passenger {
	logger.Info("MainState.CreatePassenger called")

	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	passenger, ok := s.passengers[id]
	if !ok {
		passenger = NewPassenger()
		passenger.Id = id
		passenger.CreatedAt = s.Clock.Now()
		s.passengers[id] = passenger
	}

	passenger.Active = true
	return passenger
}

// RemovePassenger marks the passenger with the given id inactive. The record is kept,
// not deleted, so it stays visible in GetPassengersSnapshot.
func (s *MainState) RemovePassenger(id string) {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	if passenger, ok := s.passengers[id]; ok {
		passenger.Active = false
	}
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
				s.CreateDriver(e.ClientId)
				fmt.Println("Driver registered:", e.ClientId)
			case events.ClientTypePassenger:
				s.CreatePassenger(e.ClientId)
				fmt.Println("Passenger registered:", e.ClientId)
			}
		case *events.ClientUnregisteredEvent:
			switch e.ClientType {
			case events.ClientTypeDriver:
				s.RemoveDriver(e.ClientId)
				fmt.Println("Driver unregistered/inactivated:", e.ClientId)
			case events.ClientTypePassenger:
				s.RemovePassenger(e.ClientId)
				fmt.Println("Passenger unregistered/inactivated:", e.ClientId)
			}
		}

	}
}

func (s *MainState) GetInvitationsContainer() *InvitationsContainer {
	return s.invitationsContainer
}
