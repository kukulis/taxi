package messages

import "darbelis.eu/taxi/internal/state"

// === to driver or to passenger

type ServerRequestsCoordinates struct{}

func (ServerRequestsCoordinates) GetMessageType() MessageType {
	return MessageTypeServerRequestsCoordinates
}

type ServerRequestsClientInfo struct{}

func (ServerRequestsClientInfo) GetMessageType() MessageType {
	return MessageTypeServerRequestsClientInfo
}

// === driver or passenger

type ClientRespondsCoordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (ClientRespondsCoordinates) GetMessageType() MessageType {
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

func (ClientRespondsInfo) GetMessageType() MessageType {
	return MessageTypeClientRespondsInfo
}

// ClientResponseError is sent by the client when it didn't understand a
// request the server sent it (e.g. an unrecognized message type).
type ClientResponseError struct {
	Error string `json:"error"`
}

func (ClientResponseError) GetMessageType() MessageType {
	return MessageTypeClientResponseError
}

// ClientError is sent by the client when it hit its own functionality
// error unrelated to a specific server request (e.g. it can't read device coordinates).
type ClientError struct {
	Error string `json:"error"`
}

func (ClientError) GetMessageType() MessageType {
	return MessageTypeClientError
}

// --- driver-related messages

type ServerOffersPassenger struct {
	PassengerId string  `json:"passenger_id"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func (ServerOffersPassenger) GetMessageType() MessageType {
	return MessageTypeServerOffersPassenger
}

type ServerCancelsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (ServerCancelsOffer) GetMessageType() MessageType {
	return MessageTypeServerCancelsOffer
}

type DriverCancelsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverCancelsOffer) GetMessageType() MessageType {
	return MessageTypeDriverCancelsOffer
}

type DriverAcceptsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverAcceptsOffer) GetMessageType() MessageType {
	return MessageTypeDriverAcceptsOffer
}

type DriverRejectsOffer struct {
	PassengerId string `json:"passenger_id"`
}

func (DriverRejectsOffer) GetMessageType() MessageType {
	return MessageTypeDriverRejectsOffer
}

type DriverChangesStatus struct {
	Status state.DriverStatus `json:"status"`
}

func (DriverChangesStatus) GetMessageType() MessageType {
	return MessageTypeDriverChangesStatus
}

// --- passenger-related messages

type PassengerInvitesDriver struct {
	DriverId string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func (PassengerInvitesDriver) GetMessageType() MessageType {
	return MessageTypePassengerInvitesDriver
}

type PassengerCancelsInvite struct {
	DriverId string `json:"driver_id"`
}

func (PassengerCancelsInvite) GetMessageType() MessageType {
	return MessageTypePassengerCancelsInvite
}

type ServerNotifiesInviteAccepted struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteAccepted) GetMessageType() MessageType {
	return MessageTypeServerNotifiesInviteAccepted
}

type ServerNotifiesInviteRejected struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteRejected) GetMessageType() MessageType {
	return MessageTypeServerNotifiesInviteRejected
}

type ServerNotifiesInviteCanceled struct {
	DriverId string `json:"driver_id"`
}

func (ServerNotifiesInviteCanceled) GetMessageType() MessageType {
	return MessageTypeServerNotifiesInviteCanceled
}

type PassengerRequestDriverCoords struct {
	DriverId string `json:"driver_id"`
}

func (PassengerRequestDriverCoords) GetMessageType() MessageType {
	return MessageTypePassengerRequestDriverCoords
}

type ServerRespondsDriverCoords struct {
	DriverId string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func (ServerRespondsDriverCoords) GetMessageType() MessageType {
	return MessageTypeServerRespondsDriverCoords
}

type ServerNotifiesVoyageStarted struct {
	DriverId    string `json:"driver_id"`
	PassengerId string `json:"passenger_id"`
}

func (ServerNotifiesVoyageStarted) GetMessageType() MessageType {
	return MessageTypeServerNotifiesVoyageStarted
}

type ServerNotifiesVoyageFinished struct {
	DriverId    string `json:"driver_id"`
	PassengerId string `json:"passenger_id"`
}

func (ServerNotifiesVoyageFinished) GetMessageType() MessageType {
	return MessageTypeServerNotifiesVoyageFinished
}
