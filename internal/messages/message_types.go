package messages

import (
	"encoding/json"

	"darbelis.eu/taxi/internal/message_common"
)

// Envelope is the wire format of every message.
// Data is decoded into a payload struct once Type is known.
type Envelope struct {
	Id   string                     `json:"id"`
	Type message_common.MessageType `json:"type"`
	Data json.RawMessage            `json:"data"`
}

const (
	// === to driver or to passenger

	MessageTypeServerRequestsCoordinates message_common.MessageType = "server_requests_coordinates"
	MessageTypeServerRequestsClientInfo  message_common.MessageType = "server_requests_client_info"

	// === driver or passenger

	MessageTypeClientRespondsCoordinates message_common.MessageType = "client_responds_coordinates"
	MessageTypeClientRespondsInfo        message_common.MessageType = "client_responds_info"
	MessageTypeClientResponseError       message_common.MessageType = "client_response_error"
	MessageTypeClientError               message_common.MessageType = "client_error"

	// --- driver-related messages

	MessageTypeServerOffersPassenger message_common.MessageType = "server_offers_passenger"
	MessageTypeServerCancelsOffer    message_common.MessageType = "server_cancels_offer"
	MessageTypeDriverCancelsOffer    message_common.MessageType = "driver_cancels_offer"
	MessageTypeDriverAcceptsOffer    message_common.MessageType = "driver_accepts_offer"
	MessageTypeDriverRejectsOffer    message_common.MessageType = "driver_rejects_offer"

	// if driver status becomes "working" and his coordinates and client coordinates are equal ( near )
	// then this will generate an event about a started voyage
	// if driver changes status from "working" to any other, then this will generate an event about a finished voyage

	MessageTypeDriverChangesStatus message_common.MessageType = "driver_changes_status"

	// --- passenger-related messages

	MessageTypePassengerInvitesDriver       message_common.MessageType = "passenger_invites_driver"
	MessageTypePassengerCancelsInvite       message_common.MessageType = "passenger_cancels_invite"
	MessageTypeServerNotifiesInviteAccepted message_common.MessageType = "server_notifies_invite_accepted"
	MessageTypeServerNotifiesInviteRejected message_common.MessageType = "server_notifies_invite_rejected"
	MessageTypeServerNotifiesInviteCanceled message_common.MessageType = "server_notifies_invite_canceled"
	MessageTypePassengerRequestDriverCoords message_common.MessageType = "passenger_request_driver_coordinates"
	MessageTypeServerRespondsDriverCoords   message_common.MessageType = "server_responds_driver_coordinates"
	MessageTypeServerNotifiesVoyageStarted  message_common.MessageType = "server_notifies_voyage_started"
	MessageTypeServerNotifiesVoyageFinished message_common.MessageType = "server_notifies_voyage_finished"

	MessageTypeServerRefreshDriverStatus         message_common.MessageType = "server_refresh_driver_status"
	MessageTypeServerRefreshDriverOffers         message_common.MessageType = "server_refresh_driver_offers"
	MessageTypeServerRefreshPassengerInvitations message_common.MessageType = "server_refresh_passenger_invitations"

	// --- virtual messages (internal notifications, not sent by any client)

	VirtualMessageTypeDriverRegistered      message_common.MessageType = "driver_registered"
	VirtualMessageTypeDriverUnregistered    message_common.MessageType = "driver_unregistered"
	VirtualMessageTypePassengerRegistered   message_common.MessageType = "passenger_registered"
	VirtualMessageTypePassengerUnregistered message_common.MessageType = "passenger_unregistered"
)
