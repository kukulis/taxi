package events

const PassengerRegisteredEventName = "PassengerRegistered"

// @deprecated use ClientRegisteredEvent instead
type PassengerRegisteredEvent struct {
	ClientId string
}

func (e *PassengerRegisteredEvent) GetName() string {
	return PassengerRegisteredEventName
}
