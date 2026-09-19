package messages

import "encoding/json"

type MessageType string

// Envelope is the wire format of every message.
// Data is decoded into a payload struct once Type is known.
type Envelope struct {
	Id   string          `json:"id"`
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MessageInterface is implemented by every message payload struct.
// The message id and the wire format are handled outside of the payload.
type MessageInterface interface {
	GetMessageType() MessageType
}

const (
	// === to driver or to passenger

	MessageTypeServerRequestsCoordinates MessageType = "server_requests_coordinates"

	// === driver or passenger

	MessageTypeClientRespondsCoordinates MessageType = "client_responds_coordinates"

	// --- driver-related messages

	MessageTypeServerOffersPassenger MessageType = "server_offers_passenger"
	MessageTypeServerCancelsOffer    MessageType = "server_cancels_offer"
	MessageTypeDriverCancelsOffer    MessageType = "driver_cancels_offer"
	MessageTypeDriverAcceptsOffer    MessageType = "driver_accepts_offer"
	MessageTypeDriverRejectsOffer    MessageType = "driver_rejects_offer"

	// if driver status becomes "working" and his coordinates and client coordinates are equal ( near )
	// then this will generate an event about a started voyage
	// if driver changes status from "working" to any other, then this will generate an event about a finished voyage

	MessageTypeDriverChangesStatus MessageType = "driver_changes_status"

	// --- passenger-related messages

	MessageTypePassengerInvitesDriver       MessageType = "passenger_invites_driver"
	MessageTypePassengerCancelsInvite       MessageType = "passenger_cancels_invite"
	MessageTypeServerNotifiesInviteAccepted MessageType = "server_notifies_invite_accepted"
	MessageTypeServerNotifiesInviteRejected MessageType = "server_notifies_invite_rejected"
	MessageTypeServerNotifiesInviteCanceled MessageType = "server_notifies_invite_canceled"
	MessageTypePassengerRequestDriverCoords MessageType = "passenger_request_driver_coordinates"
	MessageTypeServerRespondsDriverCoords   MessageType = "server_responds_driver_coordinates"
	MessageTypeServerNotifiesVoyageStarted  MessageType = "server_notifies_voyage_started"
	MessageTypeServerNotifiesVoyageFinished MessageType = "server_notifies_voyage_finished"
)
