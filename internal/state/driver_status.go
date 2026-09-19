package state

type DriverStatus string

const (
	DriverStatusIdle    DriverStatus = "idle"
	DriverStatusWorking DriverStatus = "working"
	DriverStatusOffline DriverStatus = "offline"
	DriverStatusResting DriverStatus = "resting"
)
