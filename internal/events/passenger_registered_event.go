package events

const PassengerRegisteredEventName = "PassengerRegistered"

type PassengerRegisteredEvent struct {
	ClientId string
}

func (e *PassengerRegisteredEvent) GetName() string {
	return PassengerRegisteredEventName
}
