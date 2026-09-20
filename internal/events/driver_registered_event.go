package events

const DriverRegisteredEventName = "DriverRegistered"

type DriverRegisteredEvent struct {
	ClientId string
}

func (e *DriverRegisteredEvent) GetName() string {
	return DriverRegisteredEventName
}
