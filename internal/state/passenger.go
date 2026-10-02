package state

import "time"

type Passenger struct {
	// may be a random hash for the start
	Id       string
	Lat, Lon float64
	Address  string
	Ip       string
	Phone    string
	Active   bool
	Online   bool

	CreatedAt             time.Time
	InfoReceivedAt        time.Time
	CoordinatesReceivedAt time.Time
}

func (p *Passenger) getId() string {
	return p.Id
}

func (p *Passenger) setId(id string) {
	p.Id = id
}

func (p *Passenger) copyMe() Client {
	copied := *p
	return &copied
}

func (p *Passenger) getLat() float64 {
	return p.Lat
}

func (p *Passenger) setLat(lat float64) {
	p.Lat = lat
}

func (p *Passenger) getLon() float64 {
	return p.Lon
}

func (p *Passenger) setLon(lon float64) {
	p.Lon = lon
}

func (p *Passenger) getAddress() string {
	return p.Address
}

func (p *Passenger) setAddress(address string) {
	p.Address = address
}

func (p *Passenger) getIp() string {
	return p.Ip
}

func (p *Passenger) setIp(ip string) {
	p.Ip = ip
}

func (p *Passenger) getActive() bool {
	return p.Active
}

func (p *Passenger) setActive(a bool) {
	p.Active = a
}

func (p *Passenger) getOnline() bool {
	return p.Online
}

func (p *Passenger) setOnline(online bool) {
	p.Online = online
}

func (p *Passenger) getPhone() string {
	return p.Phone
}

func (p *Passenger) setPhone(phone string) {
	p.Phone = phone
}

func (p *Passenger) getDriverStatus() DriverStatus {
	return DriverStatusUndefined
}

func (p *Passenger) setDriverStatus(status DriverStatus) {
}

func (p *Passenger) getCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Passenger) setCreatedAt(t time.Time) {
	p.CreatedAt = t
}

func (p *Passenger) getCoordinatesReceivedAt() time.Time {
	return p.CoordinatesReceivedAt
}

func (p *Passenger) setCoordinatesReceivedAt(t time.Time) {
	p.CoordinatesReceivedAt = t
}

func (p *Passenger) getType() ClientType {
	return ClientTypePassenger
}

func (p *Passenger) getKey() string {
	return ClientTypePassenger.GetKey(p.Id)
}

// ensure interface is implemented

var _ Client = &Passenger{}

func NewPassenger() *Passenger {
	return &Passenger{}
}

func GetPassengerId(p *Passenger) string { return p.Id }
