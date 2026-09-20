package events

const DriverUnregisteredEventName = "DriverUnregistered"

type DriverUnregisteredEvent struct {
	ClientId string
}

func (e *DriverUnregisteredEvent) GetName() string {
	return DriverUnregisteredEventName
}
