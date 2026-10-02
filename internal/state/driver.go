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

func (d *Driver) getId() string {
	return d.Id
}

func (d *Driver) setId(id string) {
	d.Id = id
}

func (d *Driver) copyMe() Client {
	copied := *d
	return &copied
}

func (d *Driver) getLat() float64 {
	return d.Lat
}

func (d *Driver) setLat(lat float64) {
	d.Lat = lat
}

func (d *Driver) getLon() float64 {
	return d.Lon
}

func (d *Driver) setLon(lon float64) {
	d.Lon = lon
}

func (d *Driver) getAddress() string {
	return d.Address
}

func (d *Driver) setAddress(address string) {
	d.Address = address
}

func (d *Driver) getIp() string {
	return d.Ip
}

func (d *Driver) setIp(ip string) {
	d.Ip = ip
}

func (d *Driver) getActive() bool {
	return d.Active
}

func (d *Driver) setActive(a bool) {
	d.Active = a
}

func (d *Driver) getOnline() bool {
	return d.Online
}

func (d *Driver) setOnline(online bool) {
	d.Online = online
}

func (d *Driver) getPhone() string {
	return d.Phone
}

func (d *Driver) setPhone(phone string) {
	d.Phone = phone
}

func (d *Driver) getDriverStatus() DriverStatus {
	return d.Status
}

func (d *Driver) setDriverStatus(status DriverStatus) {
	d.Status = status
}

func (d *Driver) getCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *Driver) setCreatedAt(t time.Time) {
	d.CreatedAt = t
}

func (d *Driver) getCoordinatesReceivedAt() time.Time {
	return d.CoordinatesReceivedAt
}

func (d *Driver) setCoordinatesReceivedAt(t time.Time) {
	d.CoordinatesReceivedAt = t
}

func (d *Driver) getType() ClientType {
	return ClientTypeDriver
}
func (d *Driver) getKey() string {
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
