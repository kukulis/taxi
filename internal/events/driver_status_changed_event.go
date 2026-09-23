package events

const DriverStatusChangedEventName = "DriverStatusChanged"

type DriverStatusChangedEvent struct {
	// leaving string to avoid cyclic dependencies
	DriverId        string
	DriverStatusOld string
	DriverStatusNew string
}

func (d DriverStatusChangedEvent) GetName() string {
	return DriverStatusChangedEventName
}
