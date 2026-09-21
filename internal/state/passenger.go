package state

type Passenger struct {
	// may be a random hash for the start
	Id       string
	Lat, Lon float64
	Ip       string
	Phone    string
	Active   bool
}

func NewPassenger() *Passenger {
	return &Passenger{}
}
