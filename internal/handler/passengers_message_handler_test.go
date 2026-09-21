package handler

import (
	"fmt"
	"testing"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
)

// findPassenger returns the passenger with the given id from mainState's snapshot, or nil.
func findPassenger(mainState *state.MainState, id string) *state.Passenger {
	for _, passenger := range mainState.GetPassengersSnapshot() {
		if passenger.Id == id {
			return passenger
		}
	}
	return nil
}

func TestPassengersMessageHandler_UpdatesPassengerCoordinatesOnResponse(t *testing.T) {
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.RealClock{})
	mainState.CreatePassenger(passengerId)

	// same hub mock initialization as in hub_mock_test.go: answering
	// ServerRequestsCoordinates with ClientRespondsCoordinates on the incoming channel.
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(sent ws.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsCoordinates); !ok {
			return
		}
		passengersHub.FeedIncomingMessage(ws.ClientMessage{
			ClientId: sent.ClientId,
			Message:  messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
		})
	}

	driversHub := ws.NewHubMock()

	passengersHandler := NewPassengersMessageHandler(mainState, driversHub, passengersHub)
	go passengersHandler.Handle(passengersHub.GetIncomingMessagesChannel())

	passengersHub.SendMessage(ws.ClientMessage{
		ClientId: passengerId,
		Message:  messages.ServerRequestsCoordinates{},
	})

	waitFor(t, 200*time.Millisecond,
		func() bool {
			passenger := findPassenger(mainState, passengerId)
			return passenger != nil && passenger.Lat == 54.68 && passenger.Lon == 25.28
		},
		func() string {
			passenger := findPassenger(mainState, passengerId)
			if passenger == nil {
				return fmt.Sprintf("findPassenger(%q) = nil, want passenger with Lat=54.68 Lon=25.28", passengerId)
			}
			return fmt.Sprintf("passenger coordinates = {%v %v}, want {54.68 25.28}", passenger.Lat, passenger.Lon)
		},
	)
}

func TestPassengersMessageHandler_UpdatesPassengerInfoOnResponse(t *testing.T) {
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.RealClock{})
	mainState.CreatePassenger(passengerId)

	// answers ServerRequestsClientInfo with ClientRespondsInfo on the incoming channel.
	passengersHub := ws.NewHubMock()
	passengersHub.SendMessageFunc = func(sent ws.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsClientInfo); !ok {
			return
		}
		passengersHub.FeedIncomingMessage(ws.ClientMessage{
			ClientId: sent.ClientId,
			Message:  messages.ClientRespondsInfo{Phone: "+37060054321"},
		})
	}

	driversHub := ws.NewHubMock()

	passengersHandler := NewPassengersMessageHandler(mainState, driversHub, passengersHub)
	go passengersHandler.Handle(passengersHub.GetIncomingMessagesChannel())

	passengersHub.SendMessage(ws.ClientMessage{
		ClientId: passengerId,
		Message:  messages.ServerRequestsClientInfo{},
	})

	waitFor(t, 200*time.Millisecond,
		func() bool {
			passenger := findPassenger(mainState, passengerId)
			return passenger != nil && passenger.Phone == "+37060054321"
		},
		func() string {
			passenger := findPassenger(mainState, passengerId)
			if passenger == nil {
				return fmt.Sprintf("findPassenger(%q) = nil, want passenger with Phone=+37060054321", passengerId)
			}
			return fmt.Sprintf("passenger Phone = %q, want %q", passenger.Phone, "+37060054321")
		},
	)
}
