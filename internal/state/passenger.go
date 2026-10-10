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

	Name string

	CreatedAt             time.Time
	InfoReceivedAt        time.Time
	CoordinatesReceivedAt time.Time
}

func (p *Passenger) GetId() string {
	return p.Id
}

func (p *Passenger) SetId(id string) {
	p.Id = id
}

func (p *Passenger) CopyMe() Client {
	copied := *p
	return &copied
}

func (p *Passenger) GetLat() float64 {
	return p.Lat
}

func (p *Passenger) SetLat(lat float64) {
	p.Lat = lat
}

func (p *Passenger) GetLon() float64 {
	return p.Lon
}

func (p *Passenger) SetLon(lon float64) {
	p.Lon = lon
}

func (p *Passenger) GetAddress() string {
	return p.Address
}

func (p *Passenger) SetAddress(address string) {
	p.Address = address
}

func (p *Passenger) GetIp() string {
	return p.Ip
}

func (p *Passenger) SetIp(ip string) {
	p.Ip = ip
}

func (p *Passenger) GetActive() bool {
	return p.Active
}

func (p *Passenger) SetActive(a bool) {
	p.Active = a
}

func (p *Passenger) GetOnline() bool {
	return p.Online
}

func (p *Passenger) SetOnline(online bool) {
	p.Online = online
}

func (p *Passenger) GetPhone() string {
	return p.Phone
}

func (p *Passenger) SetPhone(phone string) {
	p.Phone = phone
}

func (p *Passenger) GetDriverStatus() DriverStatus {
	return DriverStatusUndefined
}

func (p *Passenger) SetDriverStatus(status DriverStatus) {
}

func (p *Passenger) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Passenger) SetCreatedAt(t time.Time) {
	p.CreatedAt = t
}

func (p *Passenger) GetCoordinatesReceivedAt() time.Time {
	return p.CoordinatesReceivedAt
}

func (p *Passenger) SetCoordinatesReceivedAt(t time.Time) {
	p.CoordinatesReceivedAt = t
}

func (p *Passenger) GetType() ClientType {
	return ClientTypePassenger
}

func (p *Passenger) GetKey() string {
	return ClientTypePassenger.GetKey(p.Id)
}

func (p *Passenger) GetName() string {
	return p.Name
}

func (p *Passenger) SetName(name string) {
	p.Name = name
}

// ensure interface is implemented

var _ Client = &Passenger{}

func NewPassenger() *Passenger {
	return &Passenger{}
}

func GetPassengerId(p *Passenger) string { return p.Id }
