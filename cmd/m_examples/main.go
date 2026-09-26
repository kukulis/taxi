package main

import (
	"fmt"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/util"
)

// exampleMessages mirrors internal/messages/message_codec_test.go's roundTripMessages:
// one value per message type, with non-zero fields.
var exampleMessages = []message_common.MessageInterface{
	messages.ServerRequestsCoordinates{},
	messages.ServerRequestsClientInfo{},
	messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
	messages.ClientRespondsInfo{Phone: "+37060012345", VehicleInfo: "Toyota Prius", Ip: "203.0.113.5"},
	messages.ClientResponseError{Error: "unrecognized message type"},
	messages.ClientError{Error: "can't retrieve coordinates from the device"},

	messages.ServerOffersPassenger{PassengerId: "p1", Lat: 54.68, Lon: 25.28},
	messages.ServerCancelsOffer{PassengerId: "p1"},
	messages.DriverCancelsOffer{PassengerId: "p1"},
	messages.DriverAcceptsOffer{PassengerId: "p1"},
	messages.DriverRejectsOffer{PassengerId: "p1"},
	messages.DriverChangesStatus{Status: state.DriverStatusWorking},

	messages.PassengerInvitesDriver{DriverId: "d1", Lat: 54.68, Lon: 25.28},
	messages.PassengerCancelsInvite{DriverId: "d1"},
	messages.ServerNotifiesInviteAccepted{DriverId: "d1"},
	messages.ServerNotifiesInviteRejected{DriverId: "d1"},
	messages.ServerNotifiesInviteCanceled{DriverId: "d1"},
	messages.PassengerRequestDriverCoords{DriverId: "d1"},
	messages.ServerRespondsDriverCoords{DriverId: "d1", Lat: 54.68, Lon: 25.28},
	messages.ServerNotifiesVoyageStarted{DriverId: "d1", PassengerId: "p1"},
	messages.ServerNotifiesVoyageFinished{DriverId: "d1", PassengerId: "p1"},
}

func main() {
	idGenerator := &util.SimpleIdGenerator{}

	for _, msg := range exampleMessages {
		raw, err := messages.Encode(idGenerator.NextId(), msg)
		if err != nil {
			fmt.Printf("encode error for %T: %v\n", msg, err)
			continue
		}
		fmt.Println(string(raw))
	}
}
