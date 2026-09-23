package state

type DriverStatus string

const (
	DriverStatusIdle    DriverStatus = "idle"
	DriverStatusWorking DriverStatus = "working"
	DriverStatusOffline DriverStatus = "offline"
	DriverStatusResting DriverStatus = "resting"
)

func ResolveDriverStatus(status string) DriverStatus {
	// TODO validate status
	return DriverStatus(status)
}
