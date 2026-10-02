package state

import (
	"math"
	"sort"
	"sync"

	"darbelis.eu/taxi/pkg/util"
)

type MainState struct {
	// Clock provides the current time for timestamps such as Driver/Passenger
	// CreatedAt.
	Clock util.Clock

	invitationsContainer *InvitationsContainer

	clients     map[string]Client
	clientsLock sync.Mutex
}

func NewMainState(clock util.Clock) *MainState {
	return &MainState{
		//passengersLock:       sync.Mutex{},
		Clock:                clock,
		invitationsContainer: NewInvitationsContainer(),
		clients:              make(map[string]Client),
		clientsLock:          sync.Mutex{},
	}
}

type ClientDistance struct {
	Client   Client
	Distance float64
}

// RoundDistance returns a copy of cd with Distance rounded to 2 decimal
// places, so callers (e.g. test assertions) aren't exposed to float noise.
func (cd ClientDistance) RoundDistance() ClientDistance {
	cd.Distance = math.Round(cd.Distance*100) / 100
	return cd
}

// GetNearestClients searches clients in a square then calculates exact distance; returns a sorted array.
func (s *MainState) GetNearestClients(filter *ClientFilter, lat, lon, distanceKm float64) []ClientDistance {
	filter.CalculateSquare(lat, lon, distanceKm)

	clients := s.SearchClients(filter)

	result := util.ArrayMap(clients, func(c Client) ClientDistance {
		return ClientDistance{
			Client:   c,
			Distance: util.HaversineDistanceKm(lat, lon, c.getLat(), c.getLon()),
		}
	})

	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})

	return result
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

// SearchClientsByIds returns the clients of type t found for the given ids.
// Ids with no matching client are simply absent from the result.
func (s *MainState) SearchClientsByIds(ids []string, t ClientType) []Client {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	result := make([]Client, 0, len(ids))
	for _, id := range ids {
		if client, ok := s.clients[t.GetKey(id)]; ok {
			result = append(result, client.copyMe())
		}
	}

	return result
}
