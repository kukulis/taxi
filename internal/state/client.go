package state

import (
	"math"
	"time"
)

type Client interface {
	GetId() string
	SetId(id string)

	CopyMe() Client

	GetLat() float64
	SetLat(lat float64)
	GetLon() float64
	SetLon(lon float64)

	GetAddress() string
	SetAddress(address string)

	GetIp() string
	SetIp(ip string)

	GetActive() bool
	SetActive(a bool)

	GetOnline() bool
	SetOnline(online bool)

	GetPhone() string
	SetPhone(phone string)

	GetDriverStatus() DriverStatus
	SetDriverStatus(status DriverStatus)

	GetCreatedAt() time.Time
	SetCreatedAt(t time.Time)

	GetCoordinatesReceivedAt() time.Time
	SetCoordinatesReceivedAt(t time.Time)

	GetType() ClientType

	GetKey() string
}

type ClientType string

const ClientTypeDriver ClientType = "driver"
const ClientTypePassenger ClientType = "passenger"

func (t ClientType) GetKeyPrefix() string {
	if t == ClientTypePassenger {
		return "p_"
	}

	if t == ClientTypeDriver {
		return "d_"
	}

	return ""
}

func (t ClientType) GetKey(id string) string {
	return t.GetKeyPrefix() + id
}

func ResolveClientType(t string) ClientType {
	if t == string(ClientTypePassenger) {
		return ClientTypePassenger
	}

	if t == string(ClientTypeDriver) {
		return ClientTypeDriver
	}

	return ""
}

// ClientFilter used to filter clients
// empty string means ignore
type ClientFilter struct {
	Id      string
	Type    string
	Address string
	// ignore if empty slice
	DriverStatuses []DriverStatus
	VehicleInfo    string
	Phone          string
	Ip             string
	// Active -1 ignore, 0 false, 1 true
	Active int8
	// Online -1 ignore, 0 false, 1 true
	Online int8

	LatFrom float64
	LatTo   float64
	LonFrom float64
	LonTo   float64
}

func (f *ClientFilter) KeyFieldsOnly() bool {
	return f.Id != "" &&
		(f.Type == string(ClientTypePassenger) || f.Type == string(ClientTypeDriver)) &&
		f.LatFrom == 0 &&
		f.LatTo == 0 &&
		f.LonFrom == 0 &&
		f.LonTo == 0 &&
		f.Address == "" &&
		len(f.DriverStatuses) == 0 &&
		f.VehicleInfo == "" &&
		f.Phone == "" &&
		f.Ip == "" &&
		f.Active == -1 &&
		f.Online == -1
}

func NewEmptyFilter() *ClientFilter {
	return &ClientFilter{
		Active: -1,
		Online: -1,
	}
}

func (f *ClientFilter) Match(c Client) bool {
	if f.Id != "" {
		if f.Id != c.GetId() {
			return false
		}
	}

	if f.Type != "" {
		if f.Type != string(c.GetType()) {
			return false
		}
	}

	if f.LatFrom != 0 || f.LatTo != 0 {
		lat := c.GetLat()
		if lat < f.LatFrom || lat > f.LatTo {
			return false
		}
	}

	if f.LonFrom != 0 || f.LonTo != 0 {
		lon := c.GetLon()
		if lon < f.LonFrom || lon > f.LonTo {
			return false
		}
	}

	if f.Address != "" {
		if f.Address != c.GetAddress() {
			return false
		}
	}

	if f.Ip != "" {
		if f.Ip != c.GetIp() {
			return false
		}
	}

	if f.Phone != "" {
		if f.Phone != c.GetPhone() {
			return false
		}
	}

	if len(f.DriverStatuses) > 0 {
		matched := false
		for _, status := range f.DriverStatuses {
			if status == c.GetDriverStatus() {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if f.Active != -1 {
		if (f.Active == 1) != c.GetActive() {
			return false
		}
	}

	if f.Online != -1 {
		if (f.Online == 1) != c.GetOnline() {
			return false
		}
	}

	return true
}

// CalculateSquare sets LatFrom/LatTo/LonFrom/LonTo to a square centered on
// (lat, lon) whose edges are each distanceKm away from the center, along the
// great-circle approximation also used by util.HaversineDistanceKm (sphere of
// radius 6371km). Longitude degrees are scaled by cos(lat) since a degree of
// longitude covers less ground distance away from the equator.
func (f *ClientFilter) CalculateSquare(lat, lon float64, distanceKm float64) {
	const earthRadiusKm = 6371.0

	kmPerDegreeLat := earthRadiusKm * math.Pi / 180
	latDelta := distanceKm / kmPerDegreeLat
	lonDelta := distanceKm / (kmPerDegreeLat * math.Cos(lat*math.Pi/180))

	f.LatFrom = lat - latDelta
	f.LatTo = lat + latDelta
	f.LonFrom = lon - lonDelta
	f.LonTo = lon + lonDelta
}
