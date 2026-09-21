package messages

import (
	"reflect"
	"testing"

	"darbelis.eu/taxi/internal/state"
)

// one value for every message type, with non-zero fields
var roundTripMessages = []MessageInterface{
	ServerRequestsCoordinates{},
	ServerRequestsClientInfo{},
	ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
	ClientRespondsInfo{Phone: "+37060012345", VehicleInfo: "Toyota Prius"},
	ClientResponseError{Error: "unrecognized message type"},
	ClientError{Error: "can't retrieve coordinates from the device"},

	ServerOffersPassenger{PassengerId: "p1", Lat: 54.68, Lon: 25.28},
	ServerCancelsOffer{PassengerId: "p1"},
	DriverCancelsOffer{PassengerId: "p1"},
	DriverAcceptsOffer{PassengerId: "p1"},
	DriverRejectsOffer{PassengerId: "p1"},
	DriverChangesStatus{Status: state.DriverStatusWorking},

	PassengerInvitesDriver{DriverId: "d1", Lat: 54.68, Lon: 25.28},
	PassengerCancelsInvite{DriverId: "d1"},
	ServerNotifiesInviteAccepted{DriverId: "d1"},
	ServerNotifiesInviteRejected{DriverId: "d1"},
	ServerNotifiesInviteCanceled{DriverId: "d1"},
	PassengerRequestDriverCoords{DriverId: "d1"},
	ServerRespondsDriverCoords{DriverId: "d1", Lat: 54.68, Lon: 25.28},
	ServerNotifiesVoyageStarted{DriverId: "d1", PassengerId: "p1"},
	ServerNotifiesVoyageFinished{DriverId: "d1", PassengerId: "p1"},
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	for _, want := range roundTripMessages {
		t.Run(string(want.GetMessageType()), func(t *testing.T) {
			raw, err := Encode("id-1", want)
			if err != nil {
				t.Fatalf("Encode: %v", err)
			}

			id, got, err := Decode(raw)
			if err != nil {
				t.Fatalf("Decode(%s): %v", raw, err)
			}
			if id != "id-1" {
				t.Errorf("id = %q, want %q", id, "id-1")
			}
			if got.GetMessageType() != want.GetMessageType() {
				t.Errorf("type = %q, want %q", got.GetMessageType(), want.GetMessageType())
			}

			// Decode returns a pointer to the payload struct
			gotValue := reflect.ValueOf(got)
			if gotValue.Kind() != reflect.Ptr {
				t.Fatalf("Decode returned %T, want a pointer", got)
			}
			if !reflect.DeepEqual(gotValue.Elem().Interface(), want) {
				t.Errorf("payload = %+v, want %+v", gotValue.Elem().Interface(), want)
			}
		})
	}
}

// every MessageType constant must be handled by Decode
func TestDecodeHandlesAllMessageTypes(t *testing.T) {
	allTypes := []MessageType{
		MessageTypeServerRequestsCoordinates,
		MessageTypeServerRequestsClientInfo,
		MessageTypeClientRespondsCoordinates,
		MessageTypeClientRespondsInfo,
		MessageTypeClientResponseError,
		MessageTypeClientError,
		MessageTypeServerOffersPassenger,
		MessageTypeServerCancelsOffer,
		MessageTypeDriverCancelsOffer,
		MessageTypeDriverAcceptsOffer,
		MessageTypeDriverRejectsOffer,
		MessageTypeDriverChangesStatus,
		MessageTypePassengerInvitesDriver,
		MessageTypePassengerCancelsInvite,
		MessageTypeServerNotifiesInviteAccepted,
		MessageTypeServerNotifiesInviteRejected,
		MessageTypeServerNotifiesInviteCanceled,
		MessageTypePassengerRequestDriverCoords,
		MessageTypeServerRespondsDriverCoords,
		MessageTypeServerNotifiesVoyageStarted,
		MessageTypeServerNotifiesVoyageFinished,
	}

	covered := map[MessageType]bool{}
	for _, m := range roundTripMessages {
		covered[m.GetMessageType()] = true
	}
	for _, mt := range allTypes {
		if !covered[mt] {
			t.Errorf("message type %q is missing from roundTripMessages", mt)
		}
	}
}

func TestDecodeUnknownType(t *testing.T) {
	_, _, err := Decode([]byte(`{"id":"1","type":"nope","data":{}}`))
	if err == nil {
		t.Fatal("expected an error for an unknown message type")
	}
}

func TestDecodeWithoutData(t *testing.T) {
	_, msg, err := Decode([]byte(`{"id":"1","type":"server_requests_coordinates"}`))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if _, ok := msg.(*ServerRequestsCoordinates); !ok {
		t.Errorf("got %T, want *ServerRequestsCoordinates", msg)
	}
}

func TestDecodeInvalidPayload(t *testing.T) {
	_, _, err := Decode([]byte(`{"id":"1","type":"driver_accepts_offer","data":{"passenger_id":42}}`))
	if err == nil {
		t.Fatal("expected an error for a payload with a wrong field type")
	}
}
