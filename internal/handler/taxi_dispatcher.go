package handler

import (
	"log/slog"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
)

const (
	DriversCoordinatesTimeout    = time.Second * 15
	DriversInfoTimeout           = time.Hour * 12
	DriversTickFrequency         = time.Second * 15
	PassengersCoordinatesTimeout = time.Second * 60
	PassengersInfoTimeout        = time.Hour * 24
	PassengersTickFrequency      = time.Second * 15
)

type TaxiDispatcher struct {
	mainState     *state.MainState
	driversHub    ws.HubInterface
	passengersHub ws.HubInterface
	clock         util.Clock
	logger        *slog.Logger
}

func NewTaxiDispatcher(mainState *state.MainState, driversHub ws.HubInterface, passengersHub ws.HubInterface, clock util.Clock) *TaxiDispatcher {
	return &TaxiDispatcher{
		mainState:     mainState,
		driversHub:    driversHub,
		passengersHub: passengersHub,
		clock:         clock,
		logger:        slog.Default().With("handler", "TaxiDispatcher"),
	}
}

// TickForDriversUpdates dedicated to goroutine
func (t *TaxiDispatcher) TickForDriversUpdates() {
	ticker := time.NewTicker(DriversTickFrequency)
	defer ticker.Stop()

	for range ticker.C {
		//t.RequestForDriversInfos()
		t.RequestForDriversCoordinates()
	}
}

// TickForPassengersUpdates dedicated to goroutine
func (t *TaxiDispatcher) TickForPassengersUpdates() {
	ticker := time.NewTicker(PassengersTickFrequency)
	defer ticker.Stop()

	for range ticker.C {
		//t.RequestForPassengersInfos()
		t.RequestForPassengersCoordinates()
	}
}

// RequestForDriversInfos sends ServerRequestsClientInfo to every active driver whose
// info was never received, or was received longer than DriversInfoTimeout ago.
func (t *TaxiDispatcher) RequestForDriversInfos() {
	now := t.clock.Now()
	for _, driver := range t.mainState.GetDriversSnapshot() {
		if !driver.Active {
			continue
		}
		if driver.InfoReceivedAt.IsZero() || now.Sub(driver.InfoReceivedAt) > DriversInfoTimeout {
			t.driversHub.SendMessage(ws.ClientMessage{
				ClientId: driver.Id,
				Message:  messages.ServerRequestsClientInfo{},
			})
		}
	}
}

// RequestForDriversCoordinates sends ServerRequestsCoordinates to every active driver whose
// coordinates were never received, or were received longer than DriversCoordinatesTimeout ago.
func (t *TaxiDispatcher) RequestForDriversCoordinates() {
	now := t.clock.Now()
	for _, driver := range t.mainState.GetDriversSnapshot() {
		if !driver.Active {
			continue
		}
		if driver.CoordinatesReceivedAt.IsZero() || now.Sub(driver.CoordinatesReceivedAt) > DriversCoordinatesTimeout {
			t.driversHub.SendMessage(ws.ClientMessage{
				ClientId: driver.Id,
				Message:  messages.ServerRequestsCoordinates{},
			})
		}
	}
}

// RequestForPassengersInfos sends ServerRequestsClientInfo to every active passenger whose
// info was never received, or was received longer than PassengersInfoTimeout ago.
func (t *TaxiDispatcher) RequestForPassengersInfos() {
	now := t.clock.Now()
	for _, passenger := range t.mainState.GetPassengersSnapshot() {
		if !passenger.Active {
			continue
		}
		if passenger.InfoReceivedAt.IsZero() || now.Sub(passenger.InfoReceivedAt) > PassengersInfoTimeout {
			t.passengersHub.SendMessage(ws.ClientMessage{
				ClientId: passenger.Id,
				Message:  messages.ServerRequestsClientInfo{},
			})
		}
	}
}

// RequestForPassengersCoordinates sends ServerRequestsCoordinates to every active passenger whose
// coordinates were never received, or were received longer than PassengersCoordinatesTimeout ago.
func (t *TaxiDispatcher) RequestForPassengersCoordinates() {
	now := t.clock.Now()
	for _, passenger := range t.mainState.GetPassengersSnapshot() {
		if !passenger.Active {
			continue
		}
		if passenger.CoordinatesReceivedAt.IsZero() || now.Sub(passenger.CoordinatesReceivedAt) > PassengersCoordinatesTimeout {
			t.passengersHub.SendMessage(ws.ClientMessage{
				ClientId: passenger.Id,
				Message:  messages.ServerRequestsCoordinates{},
			})
		}
	}
}
