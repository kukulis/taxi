package handler

import (
	"fmt"
	"testing"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
)

// waitFor polls condition until it returns true. If timeout elapses first, it fails the
// test with the message msg produces from the state at that moment (e.g. "got X, want Y"),
// so the failure reads like a normal assertion instead of a generic timeout notice.
func waitFor(t *testing.T, timeout time.Duration, condition func() bool, msg func() string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatal(msg())
}

func TestDriversMessageHandler_UpdatesDriverCoordinatesOnResponse(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState()
	mainState.CreateDriver(driverId)

	// same drivers hub mock initialization as in hub_mock_test.go: answering
	// ServerRequestsCoordinates with ClientRespondsCoordinates on the incoming channel.
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(sent ws.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsCoordinates); !ok {
			return
		}
		driversHub.FeedIncomingMessage(ws.ClientMessage{
			ClientId: sent.ClientId,
			Message:  messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
		})
	}

	passengersHub := ws.NewHubMock()

	driversHandler := NewDriversMessageHandler(mainState, driversHub, passengersHub)
	go driversHandler.Handle(driversHub.GetIncomingMessagesChannel())

	driversHub.SendMessage(ws.ClientMessage{
		ClientId: driverId,
		Message:  messages.ServerRequestsCoordinates{},
	})

	waitFor(t, 200*time.Millisecond,
		func() bool {
			driver := mainState.GetDriverById(driverId)
			return driver != nil && driver.Lat == 54.68 && driver.Lon == 25.28
		},
		func() string {
			driver := mainState.GetDriverById(driverId)
			if driver == nil {
				return fmt.Sprintf("GetDriverById(%q) = nil, want driver with Lat=54.68 Lon=25.28", driverId)
			}
			return fmt.Sprintf("driver coordinates = {%v %v}, want {54.68 25.28}", driver.Lat, driver.Lon)
		},
	)
}

func TestDriversMessageHandler_UpdatesDriverInfoOnResponse(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState()
	mainState.CreateDriver(driverId)

	// answers ServerRequestsClientInfo with ClientRespondsInfo on the incoming channel.
	driversHub := ws.NewHubMock()
	driversHub.SendMessageFunc = func(sent ws.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsClientInfo); !ok {
			return
		}
		driversHub.FeedIncomingMessage(ws.ClientMessage{
			ClientId: sent.ClientId,
			Message:  messages.ClientRespondsInfo{Phone: "+37060012345", VehicleInfo: "Toyota Prius"},
		})
	}

	passengersHub := ws.NewHubMock()

	driversHandler := NewDriversMessageHandler(mainState, driversHub, passengersHub)
	go driversHandler.Handle(driversHub.GetIncomingMessagesChannel())

	driversHub.SendMessage(ws.ClientMessage{
		ClientId: driverId,
		Message:  messages.ServerRequestsClientInfo{},
	})

	waitFor(t, 200*time.Millisecond,
		func() bool {
			driver := mainState.GetDriverById(driverId)
			return driver != nil && driver.Phone == "+37060012345" && driver.VehicleInfo == "Toyota Prius"
		},
		func() string {
			driver := mainState.GetDriverById(driverId)
			if driver == nil {
				return fmt.Sprintf("GetDriverById(%q) = nil, want driver with Phone=+37060012345 VehicleInfo=Toyota Prius", driverId)
			}
			return fmt.Sprintf("driver info = {%q %q}, want {\"+37060012345\" \"Toyota Prius\"}", driver.Phone, driver.VehicleInfo)
		},
	)
}
