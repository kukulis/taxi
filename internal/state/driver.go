package state

type Driver struct {
	// may be a random hash for the start
	Id          string
	Lat         float64
	Lon         float64
	Status      DriverStatus
	VehicleInfo string
	Phone       string
	Ip          string
}
