package state

type DriverStatus string

const (
	DriverStatusIdle      DriverStatus = "idle"
	DriverStatusWorking   DriverStatus = "working"
	DriverStatusOffline   DriverStatus = "offline"
	DriverStatusResting   DriverStatus = "resting"
	DriverStatusUndefined DriverStatus = "undefined"
)

// ResolveDriverStatus deprecated as unused
func ResolveDriverStatus(status string) DriverStatus {
	// TODO validate status
	return DriverStatus(status)
}
