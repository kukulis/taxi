package state

import "time"

type Passenger struct {
	// may be a random hash for the start
	Id       string
	Lat, Lon float64
	Ip       string
	Phone    string
	Active   bool

	CreatedAt             time.Time
	InfoReceivedAt        time.Time
	CoordinatesReceivedAt time.Time
}

func NewPassenger() *Passenger {
	return &Passenger{}
}
