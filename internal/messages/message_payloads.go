package messages

import (
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
)

// === to driver or to passenger

type ServerRequestsCoordinates struct{}

func (ServerRequestsCoordinates) GetMessageType() message_common.MessageType {
	return MessageTypeServerRequestsCoordinates
}

type ServerRequestsClientInfo struct{}

func (ServerRequestsClientInfo) GetMessageType() message_common.MessageType {
	return MessageTypeServerRequestsClientInfo
}

// === driver or passenger

type ClientRespondsCoordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (ClientRespondsCoordinates) GetMessageType() message_common.MessageType {
	return MessageTypeClientRespondsCoordinates
}

// ClientRespondsInfo is the client's answer to ServerRequestsClientInfo.
// VehicleInfo is only meaningful for a driver client and is left empty by a passenger.
// Ip is optional; when empty, the server keeps whatever Ip it already has on record.
type ClientRespondsInfo struct {
	Phone       string `json:"phone"`
	VehicleInfo string `json:"vehicle_info,omitempty"`
	Ip          string `json:"ip,omitempty"`
}

func (ClientRespondsInfo) GetMessageType() message_common.MessageType {
	return MessageTypeClientRespondsInfo
}

// ClientResponseError is sent by the client when it didn't understand a
// request the server sent it (e.g. an unrecognized message type).
type ClientResponseError struct {
	Error string `json:"error"`
}

func (ClientResponseError) GetMessageType() message_common.MessageType {
	return MessageTypeClientResponseError
}

// ClientError is sent by the client when it hit its own functionality
// error unrelated to a specific server request (e.g. it can't read device coordinates).
type ClientError struct {
	Error string `json:"error"`
}

func (ClientError) GetMessageType() message_common.MessageType {
	return MessageTypeClientError
}

// --- driver-related messages

type ServerOffersPassenger struct {
	PassengerId string  `json:"passenger_id"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func (ServerOffersPassenger) GetMessageType() message_common.MessageType {
	return MessageTypeServerOffersPassenger
}

type ServerCancelsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (ServerCancelsOffer) GetMessageType() message_common.MessageType {
	return MessageTypeServerCancelsOffer
}

type DriverCancelsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverCancelsOffer) GetMessageType() message_common.MessageType {
	return MessageTypeDriverCancelsOffer
}

type DriverAcceptsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverAcceptsOffer) GetMessageType() message_common.MessageType {
	return MessageTypeDriverAcceptsOffer
}

type DriverRejectsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverRejectsOffer) GetMessageType() message_common.MessageType {
	return MessageTypeDriverRejectsOffer
}

type DriverChangesStatus struct {
	Status state.DriverStatus `json:"status"`
}

func (DriverChangesStatus) GetMessageType() message_common.MessageType {
	return MessageTypeDriverChangesStatus
}

// --- passenger-related messages

type PassengerInvitesDriver struct {
	DriverId string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func (PassengerInvitesDriver) GetMessageType() message_common.MessageType {
	return MessageTypePassengerInvitesDriver
}

type PassengerCancelsInvite struct {
	DriverId string `json:"driver_id"`
}

func (PassengerCancelsInvite) GetMessageType() message_common.MessageType {
	return MessageTypePassengerCancelsInvite
}

type ServerNotifiesInviteAccepted struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteAccepted) GetMessageType() message_common.MessageType {
	return MessageTypeServerNotifiesInviteAccepted
}

type ServerNotifiesInviteRejected struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteRejected) GetMessageType() message_common.MessageType {
	return MessageTypeServerNotifiesInviteRejected
}

type ServerNotifiesInviteCanceled struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteCanceled) GetMessageType() message_common.MessageType {
	return MessageTypeServerNotifiesInviteCanceled
}

type PassengerRequestDriverCoords struct {
	DriverId string `json:"driver_id"`
}

func (PassengerRequestDriverCoords) GetMessageType() message_common.MessageType {
	return MessageTypePassengerRequestDriverCoords
}

type ServerRespondsDriverCoords struct {
	DriverId string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func (ServerRespondsDriverCoords) GetMessageType() message_common.MessageType {
	return MessageTypeServerRespondsDriverCoords
}

type ServerNotifiesVoyageStarted struct {
	DriverId    string `json:"driver_id"`
	PassengerId string `json:"passenger_id"`
}

func (ServerNotifiesVoyageStarted) GetMessageType() message_common.MessageType {
	return MessageTypeServerNotifiesVoyageStarted
}

type ServerNotifiesVoyageFinished struct {
	DriverId    string `json:"driver_id"`
	PassengerId string `json:"passenger_id"`
}

func (ServerNotifiesVoyageFinished) GetMessageType() message_common.MessageType {
	return MessageTypeServerNotifiesVoyageFinished
}

// ServerRefreshDriverStatus driver id is in the envelope object
type ServerRefreshDriverStatus struct {
	Status string `json:"status"`
}

func (ServerRefreshDriverStatus) GetMessageType() message_common.MessageType {
	return MessageTypeServerRefreshDriverStatus
}

type ServerRefreshDriverOffers struct {
	Offers []OfferDto `json:"offers"`
}

func (ServerRefreshDriverOffers) GetMessageType() message_common.MessageType {
	return MessageTypeServerRefreshDriverOffers
}

type ServerRefreshPassengerInvitations struct {
	Invitations []InvitationDto `json:"invitations"`
}

func (ServerRefreshPassengerInvitations) GetMessageType() message_common.MessageType {
	return MessageTypeServerRefreshPassengerInvitations
}
