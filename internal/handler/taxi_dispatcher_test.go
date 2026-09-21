package handler

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
)

var fixedTestTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestTaxiDispatcher_RequestForDriversInfos_SendsToDriversHub(t *testing.T) {
	clock := util.NewFixedClock(fixedTestTime)
	mainState := state.NewMainState(clock)
	mainState.CreateDriver("driver-1")

	var driversSent []ws.ClientMessage
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(msg ws.ClientMessage) {
		driversSent = append(driversSent, msg)
	}

	var passengersSent []ws.ClientMessage
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws.ClientMessage) {
		passengersSent = append(passengersSent, msg)
	}

	dispatcher := NewTaxiDispatcher(mainState, driversHub, passengersHub, clock)
	dispatcher.RequestForDriversInfos()

	if len(driversSent) != 1 {
		t.Fatalf("driversHub received %d messages, want 1", len(driversSent))
	}
	if driversSent[0].ClientId != "driver-1" {
		t.Errorf("message ClientId = %q, want %q", driversSent[0].ClientId, "driver-1")
	}
	if len(passengersSent) != 0 {
		t.Errorf("passengersHub received %d messages, want 0", len(passengersSent))
	}
}

func TestTaxiDispatcher_RequestForDriversCoordinates_SendsToDriversHub(t *testing.T) {
	clock := util.NewFixedClock(fixedTestTime)
	mainState := state.NewMainState(clock)
	mainState.CreateDriver("driver-1")

	var driversSent []ws.ClientMessage
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(msg ws.ClientMessage) {
		driversSent = append(driversSent, msg)
	}

	var passengersSent []ws.ClientMessage
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws.ClientMessage) {
		passengersSent = append(passengersSent, msg)
	}

	dispatcher := NewTaxiDispatcher(mainState, driversHub, passengersHub, clock)
	dispatcher.RequestForDriversCoordinates()

	if len(driversSent) != 1 {
		t.Fatalf("driversHub received %d messages, want 1", len(driversSent))
	}
	if driversSent[0].ClientId != "driver-1" {
		t.Errorf("message ClientId = %q, want %q", driversSent[0].ClientId, "driver-1")
	}
	if len(passengersSent) != 0 {
		t.Errorf("passengersHub received %d messages, want 0", len(passengersSent))
	}
}

func TestTaxiDispatcher_RequestForPassengersInfos_SendsToPassengersHub(t *testing.T) {
	clock := util.NewFixedClock(fixedTestTime)
	mainState := state.NewMainState(clock)
	mainState.CreatePassenger("passenger-1")

	var driversSent []ws.ClientMessage
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(msg ws.ClientMessage) {
		driversSent = append(driversSent, msg)
	}

	var passengersSent []ws.ClientMessage
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws.ClientMessage) {
		passengersSent = append(passengersSent, msg)
	}

	dispatcher := NewTaxiDispatcher(mainState, driversHub, passengersHub, clock)
	dispatcher.RequestForPassengersInfos()

	if len(passengersSent) != 1 {
		t.Fatalf("passengersHub received %d messages, want 1", len(passengersSent))
	}
	if passengersSent[0].ClientId != "passenger-1" {
		t.Errorf("message ClientId = %q, want %q", passengersSent[0].ClientId, "passenger-1")
	}
	if len(driversSent) != 0 {
		t.Errorf("driversHub received %d messages, want 0", len(driversSent))
	}
}

func TestTaxiDispatcher_RequestForPassengersCoordinates_SendsToPassengersHub(t *testing.T) {
	clock := util.NewFixedClock(fixedTestTime)
	mainState := state.NewMainState(clock)
	mainState.CreatePassenger("passenger-1")

	var driversSent []ws.ClientMessage
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(msg ws.ClientMessage) {
		driversSent = append(driversSent, msg)
	}

	var passengersSent []ws.ClientMessage
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(msg ws.ClientMessage) {
		passengersSent = append(passengersSent, msg)
	}

	dispatcher := NewTaxiDispatcher(mainState, driversHub, passengersHub, clock)
	dispatcher.RequestForPassengersCoordinates()

	if len(passengersSent) != 1 {
		t.Fatalf("passengersHub received %d messages, want 1", len(passengersSent))
	}
	if passengersSent[0].ClientId != "passenger-1" {
		t.Errorf("message ClientId = %q, want %q", passengersSent[0].ClientId, "passenger-1")
	}
	if len(driversSent) != 0 {
		t.Errorf("driversHub received %d messages, want 0", len(driversSent))
	}
}
