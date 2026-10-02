package state

import (
	"sort"
	"sync"

	"darbelis.eu/taxi/pkg/util"
	"github.com/bytedance/gopkg/util/logger"
)

type MainState struct {
	drivers     map[string]*Driver
	driversLock sync.Mutex

	passengers     map[string]*Passenger
	passengersLock sync.Mutex

	// Clock provides the current time for timestamps such as Driver/Passenger
	// CreatedAt.
	Clock util.Clock

	invitationsContainer *InvitationsContainer

	clients     map[string]Client
	clientsLock sync.Mutex
}

func NewMainState(clock util.Clock) *MainState {
	return &MainState{
		drivers:              make(map[string]*Driver),
		passengers:           make(map[string]*Passenger),
		driversLock:          sync.Mutex{},
		passengersLock:       sync.Mutex{},
		Clock:                clock,
		invitationsContainer: NewInvitationsContainer(),
		clients:              make(map[string]Client),
		clientsLock:          sync.Mutex{},
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

// GetDriversByIds returns the drivers found for the given ids. Ids with no
// matching driver are simply absent from the result.
func (s *MainState) GetDriversByIds(ids []string) []*Driver {
	s.driversLock.Lock()
	defer s.driversLock.Unlock()

	driversList := make([]*Driver, 0, len(ids))
	for _, id := range ids {
		driver, ok := s.drivers[id]
		if !ok {
			continue
		}
		driverCopy := *driver
		driversList = append(driversList, &driverCopy)
	}

	return driversList
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

func (s *MainState) GetPassengerById(id string) *Passenger {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	passenger, ok := s.passengers[id]
	if !ok {
		return nil
	}

	passengerCopy := *passenger
	return &passengerCopy
}

// GetPassengersByIds returns the passengers found for the given ids. Ids with
// no matching passenger are simply absent from the result.
func (s *MainState) GetPassengersByIds(ids []string) []*Passenger {
	s.passengersLock.Lock()
	defer s.passengersLock.Unlock()

	passengersList := make([]*Passenger, 0, len(ids))
	for _, id := range ids {
		passenger, ok := s.passengers[id]
		if !ok {
			continue
		}
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

func (s *MainState) GetInvitationsContainer() *InvitationsContainer {
	return s.invitationsContainer
}

func (s *MainState) AddClient(c Client) Client {
	if c == nil {
		return nil
	}
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	existing, ok := s.clients[c.getKey()]

	if !ok {
		s.clients[c.getKey()] = c
		existing = c
	}

	// we will differentiate these flags later
	existing.setOnline(true)
	existing.setActive(true)

	return c
}

func (s *MainState) UpdateClient(c Client, mutate func(Client)) bool {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	existing, ok := s.clients[c.getKey()]

	if !ok {
		return false
	}

	mutate(existing)

	return true
}

// RemoveClient removes client from the list
// Use UpdateClient if you want only inactivate it
func (s *MainState) RemoveClient(id string, t ClientType) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	delete(s.clients, t.GetKey(id))
}

func (s *MainState) SearchClients(f *ClientFilter) []Client {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	var result []Client

	if f.KeyFieldsOnly() {
		key := ResolveClientType(f.Type).GetKey(f.Id)

		existing, ok := s.clients[key]

		if !ok {
			return result
		}

		result = append(result, existing.copyMe())

		return result
	}

	for _, client := range s.clients {
		if f.Match(client) {
			result = append(result, client.copyMe())
		}
	}

	return result
}
