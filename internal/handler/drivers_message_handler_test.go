package handler

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/util"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

// findDriver returns the driver with the given id, or nil.
// TODO consider to make such a function as a method for the MainState
func findDriver(mainState *state.MainState, id string) *state.Driver {
	filter := state.NewEmptyFilter()
	filter.Id = id
	filter.Type = string(state.ClientTypeDriver)
	found := mainState.SearchClients(filter)
	if len(found) == 0 {
		return nil
	}
	return found[0].(*state.Driver)
}

func TestDriversMessageHandler_UpdatesDriverCoordinatesOnResponse(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	// same drivers hub mock initialization as in hub_mock_test.go: answering
	// ServerRequestsCoordinates with ClientRespondsCoordinates on the incoming channel.
	driversHubMock := ws2.NewHubMock()
	driversHubMock.SendMessageFunc = func(sent ws2.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsCoordinates); !ok {
			return
		}
		driversHubMock.FeedIncomingMessage(ws2.ClientMessage{
			ClientId: sent.ClientId,
			Message:  &messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
		})
	}

	passengersHub := ws2.NewHubMock()

	driversHandler := NewDriversMessageHandler(mainState, driversHubMock, passengersHub, &util.Dispatcher{}, nil)

	driversHubMock.SendMessage(ws2.ClientMessage{
		ClientId: driverId,
		Message:  messages.ServerRequestsCoordinates{},
	})
	driversHubMock.Close()

	driversHandler.Handle(driversHubMock.GetIncomingMessagesChannel())

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
		t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28", driver)
	}
}

// TestDriversMessageHandler_SurvivesErrorOnCoordinatesRequest feeds the handler the case's
// error reply, then a valid ClientRespondsCoordinates. It asserts the coordinates end up in
// MainState, which only happens if the handler survives the error message and keeps consuming
// the channel instead of getting stuck on it.
func TestDriversMessageHandler_SurvivesErrorOnCoordinatesRequest(t *testing.T) {
	cases := []struct {
		name       string
		errorReply message_common.MessageInterface
	}{
		{
			name:       "ClientResponseError",
			errorReply: &messages.ClientResponseError{Error: "didn't understand ServerRequestsCoordinates"},
		},
		{
			name:       "ClientError",
			errorReply: &messages.ClientError{Error: "can't retrieve coordinates from the device"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			const driverId = "driver-1"

			mainState := state.NewMainState(util.NewFixedClock(time.Now()))
			mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

			driversHubMock := ws2.NewHubMock()

			passengersHub := ws2.NewHubMock()

			driversHandler := NewDriversMessageHandler(mainState, driversHubMock, passengersHub, &util.Dispatcher{}, nil)

			driversHubMock.FeedIncomingMessage(ws2.ClientMessage{ClientId: driverId, Message: tc.errorReply})
			driversHubMock.FeedIncomingMessage(ws2.ClientMessage{ClientId: driverId, Message: &messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28}})

			driversHubMock.Close()

			driversHandler.Handle(driversHubMock.GetIncomingMessagesChannel())

			driver := findDriver(mainState, driverId)
			if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
				t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28 after %T then a valid response (handler may have gotten stuck on the error message)", driver, tc.errorReply)
			}
		})
	}
}

// TestDriversMessageHandler_HandlesRealDecodedMessage reproduces a real client's message as it
// actually arrives at the server: encoded to wire bytes, then decoded via messages.Decode, which
// always returns a pointer (see message_codec.go's Decode). The other tests in this file instead
// feed the mock a value literal (messages.ClientRespondsCoordinates{...}), which is not what a
// real client's message looks like once decoded, so they don't catch a switch that only matches
// value types and never the pointer types Decode produces.
func TestDriversMessageHandler_HandlesRealDecodedMessage(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	driversHubMock := ws2.NewHubMock()
	passengersHubMock := ws2.NewHubMock()

	driversHandler := NewDriversMessageHandler(mainState, driversHubMock, passengersHubMock, &util.Dispatcher{}, nil)

	raw, err := messages.Encode("id-1", messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28})
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	_, decoded, err := messages.Decode(raw)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	driversHubMock.FeedIncomingMessage(ws2.ClientMessage{
		ClientId: driverId,
		Message:  decoded,
	})
	driversHubMock.Close()

	driversHandler.Handle(driversHubMock.GetIncomingMessagesChannel())

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
		t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28 — Handle's type switch may not match %T, the pointer type Decode actually produces", driver, decoded)
	}
}

// Wires the handler to the real listener, so the event the handler dispatches must be the
// type the listener asserts (a pointer/value mismatch here panicked the server).
func TestDriversMessageHandler_ChangesStatusThroughListener(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	dispatcher := util.NewDispatcher()
	dispatcher.AddListener(events.DriverStatusChangedEventName, NewDriverStatusChangeForDriverListener(mainState).Handle)

	driversHub := ws2.NewHubMock()
	driversHandler := NewDriversMessageHandler(mainState, driversHub, ws2.NewHubMock(), dispatcher, nil)

	driversHub.FeedIncomingMessage(ws2.ClientMessage{
		ClientId: driverId,
		Message:  &messages.DriverChangesStatus{Status: state.DriverStatusIdle},
	})
	driversHub.Close()

	driversHandler.Handle(driversHub.GetIncomingMessagesChannel())

	if got := findDriver(mainState, driverId).Status; got != state.DriverStatusIdle {
		t.Errorf("driver.Status = %q, want %q", got, state.DriverStatusIdle)
	}
}

func TestDriversMessageHandler_UpdatesDriverInfoOnResponse(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	// answers ServerRequestsClientInfo with ClientRespondsInfo on the incoming channel.
	driversHub := ws2.NewHubMock()
	driversHub.SendMessageFunc = func(sent ws2.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsClientInfo); !ok {
			return
		}
		driversHub.FeedIncomingMessage(ws2.ClientMessage{
			ClientId: sent.ClientId,
			Message:  &messages.ClientRespondsInfo{Phone: "+37060012345", VehicleInfo: "Toyota Prius"},
		})
	}

	passengersHub := ws2.NewHubMock()

	driversHandler := NewDriversMessageHandler(mainState, driversHub, passengersHub, &util.Dispatcher{}, nil)

	driversHub.SendMessage(ws2.ClientMessage{
		ClientId: driverId,
		Message:  messages.ServerRequestsClientInfo{},
	})
	driversHub.Close()

	driversHandler.Handle(driversHub.GetIncomingMessagesChannel())

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Phone != "+37060012345" || driver.VehicleInfo != "Toyota Prius" {
		t.Errorf("driver = %+v, want driver with Phone=+37060012345 VehicleInfo=Toyota Prius", driver)
	}
}
