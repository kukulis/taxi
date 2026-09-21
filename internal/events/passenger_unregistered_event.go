package events

const PassengerUnregisteredEventName = "PassengerUnregistered"

// PassengerUnregisteredEvent @Deprecated use ClientUnregisteredEvent instead
type PassengerUnregisteredEvent struct {
	ClientId string
}

func (e *PassengerUnregisteredEvent) GetName() string {
	return PassengerUnregisteredEventName
}
