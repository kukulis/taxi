package messages

import (
	"encoding/json"
	"fmt"

	"darbelis.eu/taxi/internal/message_common"
)

// Encode wraps the payload into an Envelope and marshals it.
// The type is taken from the payload, so it can not disagree with the data.
func Encode(id string, msg message_common.MessageInterface) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return json.Marshal(Envelope{
		Id:   id,
		Type: msg.GetMessageType(),
		Data: data,
	})
}

// Decode unmarshals the Envelope, then the payload into the struct matching the type.
// The returned message is a pointer to the payload struct.
func Decode(raw []byte) (string, message_common.MessageInterface, error) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return "", nil, err
	}

	var msg message_common.MessageInterface
	switch env.Type {
	case MessageTypeServerRequestsCoordinates:
		msg = &ServerRequestsCoordinates{}
	case MessageTypeServerRequestsClientInfo:
		msg = &ServerRequestsClientInfo{}
	case MessageTypeClientRespondsCoordinates:
		msg = &ClientRespondsCoordinates{}
	case MessageTypeClientRespondsInfo:
		msg = &ClientRespondsInfo{}
	case MessageTypeClientResponseError:
		msg = &ClientResponseError{}
	case MessageTypeClientError:
		msg = &ClientError{}

	case MessageTypeServerOffersPassenger:
		msg = &ServerOffersPassenger{}
	case MessageTypeServerCancelsOffer:
		msg = &ServerCancelsOffer{}
	case MessageTypeDriverCancelsOffer:
		msg = &DriverCancelsOffer{}
	case MessageTypeDriverAcceptsOffer:
		msg = &DriverAcceptsOffer{}
	case MessageTypeDriverRejectsOffer:
		msg = &DriverRejectsOffer{}
	case MessageTypeDriverChangesStatus:
		msg = &DriverChangesStatus{}

	case MessageTypePassengerInvitesDriver:
		msg = &PassengerInvitesDriver{}
	case MessageTypePassengerCancelsInvite:
		msg = &PassengerCancelsInvite{}
	case MessageTypeServerNotifiesInviteAccepted:
		msg = &ServerNotifiesInviteAccepted{}
	case MessageTypeServerNotifiesInviteRejected:
		msg = &ServerNotifiesInviteRejected{}
	case MessageTypeServerNotifiesInviteCanceled:
		msg = &ServerNotifiesInviteCanceled{}
	case MessageTypePassengerRequestDriverCoords:
		msg = &PassengerRequestDriverCoords{}
	case MessageTypeServerRespondsDriverCoords:
		msg = &ServerRespondsDriverCoords{}
	case MessageTypeServerNotifiesVoyageStarted:
		msg = &ServerNotifiesVoyageStarted{}
	case MessageTypeServerNotifiesVoyageFinished:
		msg = &ServerNotifiesVoyageFinished{}
	case MessageTypeServerRefreshDriverStatus:
		msg = &ServerRefreshDriverStatus{}
	case MessageTypeServerRefreshDriverOffers:
		msg = &ServerRefreshDriverOffers{}
	case MessageTypeServerRefreshPassengerInvitations:
		msg = &ServerRefreshPassengerInvitations{}

	case VirtualMessageTypeDriverRegistered:
		msg = &DriverRegistered{}
	case VirtualMessageTypeDriverUnregistered:
		msg = &DriverUnregistered{}
	case VirtualMessageTypePassengerRegistered:
		msg = &PassengerRegistered{}
	case VirtualMessageTypePassengerUnregistered:
		msg = &PassengerUnregistered{}

	default:
		return env.Id, nil, fmt.Errorf("unknown message type %q", env.Type)
	}

	// payloads without fields may come without data
	if len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, msg); err != nil {
			return env.Id, nil, fmt.Errorf("decoding %q payload: %w", env.Type, err)
		}
	}

	return env.Id, msg, nil
}
