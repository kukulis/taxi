package events

const PassengerUnregisteredEventName = "PassengerUnregistered"

type PassengerUnregisteredEvent struct {
	ClientId string
}

func (e *PassengerUnregisteredEvent) GetName() string {
	return PassengerUnregisteredEventName
}
