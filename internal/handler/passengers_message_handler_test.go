package handler

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

// findPassenger returns the passenger with the given id, or nil.
func findPassenger(mainState *state.MainState, id string) *state.Passenger {
	filter := state.NewEmptyFilter()
	filter.Id = id
	filter.Type = string(state.ClientTypePassenger)
	found := mainState.SearchClients(filter)
	if len(found) == 0 {
		return nil
	}
	return found[0].(*state.Passenger)
}

func TestPassengersMessageHandler_UpdatesPassengerCoordinatesOnResponse(t *testing.T) {
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Passenger{Id: passengerId, CreatedAt: mainState.Clock.Now()})

	// same hub mock initialization as in hub_mock_test.go: answering
	// ServerRequestsCoordinates with ClientRespondsCoordinates on the incoming channel.
	passengersHub := ws2.NewHubMock()
	passengersHub.SendMessageFunc = func(sent ws2.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsCoordinates); !ok {
			return
		}
		passengersHub.FeedIncomingMessage(ws2.ClientMessage{
			ClientId: sent.ClientId,
			Message:  &messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
		})
	}

	driversHub := ws2.NewHubMock()

	passengersHandler := NewPassengersMessageHandler(mainState, driversHub, passengersHub, nil)

	passengersHub.SendMessage(ws2.ClientMessage{
		ClientId: passengerId,
		Message:  messages.ServerRequestsCoordinates{},
	})
	passengersHub.Close()

	passengersHandler.Handle(passengersHub.GetIncomingMessagesChannel())

	passenger := findPassenger(mainState, passengerId)
	if passenger == nil || passenger.Lat != 54.68 || passenger.Lon != 25.28 {
		t.Errorf("passenger = %+v, want passenger with Lat=54.68 Lon=25.28", passenger)
	}
}

func TestPassengersMessageHandler_UpdatesPassengerInfoOnResponse(t *testing.T) {
	const passengerId = "passenger-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Passenger{Id: passengerId, CreatedAt: mainState.Clock.Now()})

	// answers ServerRequestsClientInfo with ClientRespondsInfo on the incoming channel.
	passengersHub := ws2.NewHubMock()
	passengersHub.SendMessageFunc = func(sent ws2.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsClientInfo); !ok {
			return
		}
		passengersHub.FeedIncomingMessage(ws2.ClientMessage{
			ClientId: sent.ClientId,
			Message:  &messages.ClientRespondsInfo{Phone: "+37060054321"},
		})
	}

	driversHub := ws2.NewHubMock()

	passengersHandler := NewPassengersMessageHandler(mainState, driversHub, passengersHub, nil)

	passengersHub.SendMessage(ws2.ClientMessage{
		ClientId: passengerId,
		Message:  messages.ServerRequestsClientInfo{},
	})
	passengersHub.Close()

	passengersHandler.Handle(passengersHub.GetIncomingMessagesChannel())

	passenger := findPassenger(mainState, passengerId)
	if passenger == nil || passenger.Phone != "+37060054321" {
		t.Errorf("passenger = %+v, want passenger with Phone=+37060054321", passenger)
	}
}
