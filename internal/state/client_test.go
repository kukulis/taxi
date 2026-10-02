package state

import "testing"

func TestClientFilter_Match(t *testing.T) {
	cases := []struct {
		name   string
		client Client
		filter *ClientFilter
		want   bool
	}{
		{
			name:   "empty filter matches driver",
			client: &Driver{Id: "d1", Lat: 1, Lon: 2, Address: "addr", Ip: "1.2.3.4", Phone: "+370", Status: DriverStatusIdle, Active: true, Online: true},
			filter: NewEmptyFilter(),
			want:   true,
		},
		{
			name:   "empty filter matches passenger",
			client: &Passenger{Id: "p1"},
			filter: NewEmptyFilter(),
			want:   true,
		},
		{
			name:   "id match",
			client: &Driver{Id: "d1"},
			filter: &ClientFilter{Id: "d1", Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "id mismatch",
			client: &Driver{Id: "d1"},
			filter: &ClientFilter{Id: "d2", Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "type match driver",
			client: &Driver{},
			filter: &ClientFilter{Type: string(ClientTypeDriver), Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "type mismatch for passenger",
			client: &Passenger{},
			filter: &ClientFilter{Type: string(ClientTypeDriver), Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "lat match",
			client: &Driver{Lat: 54.68},
			filter: &ClientFilter{Lat: 54.68, Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "lat mismatch",
			client: &Driver{Lat: 54.68},
			filter: &ClientFilter{Lat: 10, Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "lon match",
			client: &Driver{Lon: 25.28},
			filter: &ClientFilter{Lon: 25.28, Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "lon mismatch",
			client: &Driver{Lon: 25.28},
			filter: &ClientFilter{Lon: 10, Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "address match",
			client: &Driver{Address: "Kaunas"},
			filter: &ClientFilter{Address: "Kaunas", Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "address mismatch",
			client: &Driver{Address: "Kaunas"},
			filter: &ClientFilter{Address: "Vilnius", Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "ip match",
			client: &Driver{Ip: "1.2.3.4"},
			filter: &ClientFilter{Ip: "1.2.3.4", Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "ip mismatch",
			client: &Driver{Ip: "1.2.3.4"},
			filter: &ClientFilter{Ip: "5.6.7.8", Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "phone match",
			client: &Driver{Phone: "+37060000000"},
			filter: &ClientFilter{Phone: "+37060000000", Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "phone mismatch",
			client: &Driver{Phone: "+37060000000"},
			filter: &ClientFilter{Phone: "+37069999999", Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "driver statuses: matches one of several",
			client: &Driver{Status: DriverStatusIdle},
			filter: &ClientFilter{DriverStatuses: []DriverStatus{DriverStatusIdle, DriverStatusWorking}, Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "driver statuses: matches none",
			client: &Driver{Status: DriverStatusOffline},
			filter: &ClientFilter{DriverStatuses: []DriverStatus{DriverStatusIdle, DriverStatusWorking}, Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "driver statuses: empty slice is ignored",
			client: &Driver{Status: DriverStatusOffline},
			filter: &ClientFilter{Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "driver statuses: passenger matches DriverStatusUndefined",
			client: &Passenger{},
			filter: &ClientFilter{DriverStatuses: []DriverStatus{DriverStatusUndefined}, Active: -1, Online: -1},
			want:   true,
		},
		{
			name:   "driver statuses: passenger does not match a defined status",
			client: &Passenger{},
			filter: &ClientFilter{DriverStatuses: []DriverStatus{DriverStatusIdle}, Active: -1, Online: -1},
			want:   false,
		},
		{
			name:   "active required true and matches",
			client: &Driver{Active: true},
			filter: &ClientFilter{Active: 1, Online: -1},
			want:   true,
		},
		{
			name:   "active required true but client is inactive",
			client: &Driver{Active: false},
			filter: &ClientFilter{Active: 1, Online: -1},
			want:   false,
		},
		{
			name:   "active required false and matches",
			client: &Driver{Active: false},
			filter: &ClientFilter{Active: 0, Online: -1},
			want:   true,
		},
		{
			name:   "online required true and matches",
			client: &Driver{Online: true},
			filter: &ClientFilter{Active: -1, Online: 1},
			want:   true,
		},
		{
			name:   "online required true but client is offline",
			client: &Driver{Online: false},
			filter: &ClientFilter{Active: -1, Online: 1},
			want:   false,
		},
		{
			name:   "online required false and matches",
			client: &Driver{Online: false},
			filter: &ClientFilter{Active: -1, Online: 0},
			want:   true,
		},
		{
			name:   "combined filter fails on the second checked field",
			client: &Driver{Id: "d1", Active: true},
			filter: &ClientFilter{Id: "d1", Active: 0, Online: -1},
			want:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.filter.Match(c.client)
			if got != c.want {
				t.Errorf("Match() = %v, want %v", got, c.want)
			}
		})
	}
}
