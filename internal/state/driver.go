package state

import "time"

type Driver struct {
	// may be a random hash for the start
	Id          string
	Lat         float64
	Lon         float64
	Status      DriverStatus
	VehicleInfo string
	Phone       string
	Ip          string
	Active      bool

	CreatedAt             time.Time
	InfoReceivedAt        time.Time
	CoordinatesReceivedAt time.Time
}

func NewDriver() *Driver {
	return &Driver{}
}
