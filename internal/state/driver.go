package state

import "time"

type Driver struct {
	// may be a random hash for the start
	Id          string
	Lat         float64
	Lon         float64
	Address     string
	Status      DriverStatus
	VehicleInfo string
	Phone       string
	Ip          string
	Active      bool
	Online      bool

	CreatedAt             time.Time
	InfoReceivedAt        time.Time
	CoordinatesReceivedAt time.Time
}

func (d *Driver) GetId() string {
	return d.Id
}

func (d *Driver) SetId(id string) {
	d.Id = id
}

func (d *Driver) CopyMe() Client {
	copied := *d
	return &copied
}

func (d *Driver) GetLat() float64 {
	return d.Lat
}

func (d *Driver) SetLat(lat float64) {
	d.Lat = lat
}

func (d *Driver) GetLon() float64 {
	return d.Lon
}

func (d *Driver) SetLon(lon float64) {
	d.Lon = lon
}

func (d *Driver) GetAddress() string {
	return d.Address
}

func (d *Driver) SetAddress(address string) {
	d.Address = address
}

func (d *Driver) GetIp() string {
	return d.Ip
}

func (d *Driver) SetIp(ip string) {
	d.Ip = ip
}

func (d *Driver) GetActive() bool {
	return d.Active
}

func (d *Driver) SetActive(a bool) {
	d.Active = a
}

func (d *Driver) GetOnline() bool {
	return d.Online
}

func (d *Driver) SetOnline(online bool) {
	d.Online = online
}

func (d *Driver) GetPhone() string {
	return d.Phone
}

func (d *Driver) SetPhone(phone string) {
	d.Phone = phone
}

func (d *Driver) GetDriverStatus() DriverStatus {
	return d.Status
}

func (d *Driver) SetDriverStatus(status DriverStatus) {
	d.Status = status
}

func (d *Driver) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *Driver) SetCreatedAt(t time.Time) {
	d.CreatedAt = t
}

func (d *Driver) GetCoordinatesReceivedAt() time.Time {
	return d.CoordinatesReceivedAt
}

func (d *Driver) SetCoordinatesReceivedAt(t time.Time) {
	d.CoordinatesReceivedAt = t
}

func (d *Driver) GetType() ClientType {
	return ClientTypeDriver
}
func (d *Driver) GetKey() string {
	return ClientTypeDriver.GetKey(d.Id)
}

// ensure interface is implemented

var _ Client = &Driver{}

func NewDriver() *Driver {
	return &Driver{}
}

func GetDriverId(d *Driver) string {
	return d.Id
}
