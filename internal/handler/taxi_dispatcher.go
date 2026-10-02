package handler

import (
	"log/slog"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

const (
	DriversCoordinatesTimeout = time.Second * 900
	DriversInfoTimeout        = time.Hour * 12
	DriversTickFrequency      = time.Second * 15

	PassengersCoordinatesTimeout = time.Second * 900
	PassengersInfoTimeout        = time.Hour * 24
	PassengersTickFrequency      = time.Second * 15
)

type TaxiDispatcher struct {
	mainState     *state.MainState
	driversHub    ws2.HubInterface
	passengersHub ws2.HubInterface
	clock         util.Clock
	logger        *slog.Logger
}

func NewTaxiDispatcher(mainState *state.MainState, driversHub ws2.HubInterface, passengersHub ws2.HubInterface, clock util.Clock) *TaxiDispatcher {
	return &TaxiDispatcher{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
		clock:         clock,
		logger:        slog.Default().With("handler", "TaxiDispatcher"),
	}
}

// TickForDriversUpdates dedicated to goroutine
func (tDispatcher *TaxiDispatcher) TickForDriversUpdates() {
	ticker := time.NewTicker(DriversTickFrequency)
	defer ticker.Stop()

	for range ticker.C {
		//t.RequestForDriversInfos()
		tDispatcher.RequestForDriversCoordinates()
	}
}

// TickForPassengersUpdates dedicated to goroutine
func (tDispatcher *TaxiDispatcher) TickForPassengersUpdates() {
	ticker := time.NewTicker(PassengersTickFrequency)
	defer ticker.Stop()

	for range ticker.C {
		//t.RequestForPassengersInfos()
		tDispatcher.RequestForPassengersCoordinates()
	}
}

// RequestForDriversInfos sends ServerRequestsClientInfo to every active driver whose
// info was never received, or was received longer than DriversInfoTimeout ago.
func (tDispatcher *TaxiDispatcher) RequestForDriversInfos() {
	now := tDispatcher.clock.Now()
	for _, driver := range tDispatcher.driversSnapshot() {
		if !driver.Active {
			continue
		}
		if driver.InfoReceivedAt.IsZero() || now.Sub(driver.InfoReceivedAt) > DriversInfoTimeout {
			tDispatcher.driversHub.SendMessage(ws2.ClientMessage{
				ClientId: driver.Id,
				Message:  messages.ServerRequestsClientInfo{},
			})
		}
	}
}

// RequestForDriversCoordinates sends ServerRequestsCoordinates to every active driver whose
// coordinates were never received, or were received longer than DriversCoordinatesTimeout ago.
func (tDispatcher *TaxiDispatcher) RequestForDriversCoordinates() {
	now := tDispatcher.clock.Now()
	for _, driver := range tDispatcher.driversSnapshot() {
		if !driver.Active {
			continue
		}
		if driver.CoordinatesReceivedAt.IsZero() || now.Sub(driver.CoordinatesReceivedAt) > DriversCoordinatesTimeout {
			tDispatcher.driversHub.SendMessage(ws2.ClientMessage{
				ClientId: driver.Id,
				Message:  messages.ServerRequestsCoordinates{},
			})
		}
	}
}

// RequestForPassengersInfos sends ServerRequestsClientInfo to every active passenger whose
// info was never received, or was received longer than PassengersInfoTimeout ago.
func (tDispatcher *TaxiDispatcher) RequestForPassengersInfos() {
	now := tDispatcher.clock.Now()
	for _, passenger := range tDispatcher.passengersSnapshot() {
		if !passenger.Active {
			continue
		}
		if passenger.InfoReceivedAt.IsZero() || now.Sub(passenger.InfoReceivedAt) > PassengersInfoTimeout {
			tDispatcher.passengersHub.SendMessage(ws2.ClientMessage{
				ClientId: passenger.Id,
				Message:  messages.ServerRequestsClientInfo{},
			})
		}
	}
}

// RequestForPassengersCoordinates sends ServerRequestsCoordinates to every active passenger whose
// coordinates were never received, or were received longer than PassengersCoordinatesTimeout ago.
func (tDispatcher *TaxiDispatcher) RequestForPassengersCoordinates() {
	now := tDispatcher.clock.Now()
	for _, passenger := range tDispatcher.passengersSnapshot() {
		if !passenger.Active {
			continue
		}
		if passenger.CoordinatesReceivedAt.IsZero() || now.Sub(passenger.CoordinatesReceivedAt) > PassengersCoordinatesTimeout {
			tDispatcher.passengersHub.SendMessage(ws2.ClientMessage{
				ClientId: passenger.Id,
				Message:  messages.ServerRequestsCoordinates{},
			})
		}
	}
}

func (tDispatcher *TaxiDispatcher) driversSnapshot() []*state.Driver {
	filter := state.NewEmptyFilter()
	filter.Type = string(state.ClientTypeDriver)
	filter.DriverStatuses = []state.DriverStatus{state.DriverStatusIdle, state.DriverStatusWorking}
	//filter.Online = 1
	clients := tDispatcher.mainState.SearchClients(filter)
	return util.ArrayMap(clients, func(c state.Client) *state.Driver { return c.(*state.Driver) })
}

func (tDispatcher *TaxiDispatcher) passengersSnapshot() []*state.Passenger {
	filter := state.NewEmptyFilter()
	filter.Type = string(state.ClientTypePassenger)
	clients := tDispatcher.mainState.SearchClients(filter)
	return util.ArrayMap(clients, func(c state.Client) *state.Passenger { return c.(*state.Passenger) })
}
