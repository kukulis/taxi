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

	driversHandler := NewDriversMessageHandler(mainState, ws2.NewHubMock(), ws2.NewHubMock(), &util.Dispatcher{}, nil)

	driversHandler.HandleClientRespondsCoordinates(driverId, &messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28})

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
		t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28", driver)
	}
}

// TestDriversMessageHandler_SurvivesErrorOnCoordinatesRequest calls the handler with the case's
// error reply, then with a valid ClientRespondsCoordinates. It asserts the coordinates end up in
// MainState, which only happens if handling the error reply left the handler usable for the next call.
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

			driversHandler := NewDriversMessageHandler(mainState, ws2.NewHubMock(), ws2.NewHubMock(), &util.Dispatcher{}, nil)

			switch tc.errorReply.(type) {
			case *messages.ClientResponseError:
				driversHandler.HandleClientResponseError(driverId, tc.errorReply)
			case *messages.ClientError:
				driversHandler.HandleClientError(driverId, tc.errorReply)
			}

			driversHandler.HandleClientRespondsCoordinates(driverId, &messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28})

			driver := findDriver(mainState, driverId)
			if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
				t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28 after %T then a valid response (handler may be left in a bad state)", driver, tc.errorReply)
			}
		})
	}
}

// TestDriversMessageHandler_HandlesRealDecodedMessage reproduces a real client's message as it
// actually arrives at the server: encoded to wire bytes, then decoded via messages.Decode, which
// always returns a pointer (see message_codec.go's Decode). This catches a handler whose type
// assertion only matches a value type and never the pointer type Decode actually produces.
func TestDriversMessageHandler_HandlesRealDecodedMessage(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	driversHandler := NewDriversMessageHandler(mainState, ws2.NewHubMock(), ws2.NewHubMock(), &util.Dispatcher{}, nil)

	raw, err := messages.Encode("id-1", messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28})
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	_, decoded, err := messages.Decode(raw)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	driversHandler.HandleClientRespondsCoordinates(driverId, decoded)

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Lat != 54.68 || driver.Lon != 25.28 {
		t.Errorf("driver = %+v, want driver with Lat=54.68 Lon=25.28 — the handler's type assertion may not match %T, the pointer type Decode actually produces", driver, decoded)
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

	driversHandler := NewDriversMessageHandler(mainState, ws2.NewHubMock(), ws2.NewHubMock(), dispatcher, nil)

	driversHandler.HandleDriverChangesStatus(driverId, &messages.DriverChangesStatus{Status: state.DriverStatusIdle})

	if got := findDriver(mainState, driverId).Status; got != state.DriverStatusIdle {
		t.Errorf("driver.Status = %q, want %q", got, state.DriverStatusIdle)
	}
}

func TestDriversMessageHandler_UpdatesDriverInfoOnResponse(t *testing.T) {
	const driverId = "driver-1"

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: driverId, CreatedAt: mainState.Clock.Now()})

	driversHandler := NewDriversMessageHandler(mainState, ws2.NewHubMock(), ws2.NewHubMock(), &util.Dispatcher{}, nil)

	driversHandler.HandleClientRespondsInfo(driverId, &messages.ClientRespondsInfo{Phone: "+37060012345", VehicleInfo: "Toyota Prius"})

	driver := findDriver(mainState, driverId)
	if driver == nil || driver.Phone != "+37060012345" || driver.VehicleInfo != "Toyota Prius" {
		t.Errorf("driver = %+v, want driver with Phone=+37060012345 VehicleInfo=Toyota Prius", driver)
	}
}
