package state

import (
	"sort"
	"testing"
	"time"

	"darbelis.eu/taxi/pkg/test_utils"
	"darbelis.eu/taxi/pkg/util"
)

// All per-field filter checks are covered by TestClientFilter_Match. This test
// only exercises SearchClients' own branching (key-fields-only lookup vs. the
// full scan), so each case picks a filter shape rather than individual fields.
func TestMainState_SearchClients(t *testing.T) {
	s := NewMainState(util.NewFixedClock(time.Now()))

	s.AddClient(&Driver{Id: "d1"})
	s.AddClient(&Driver{Id: "d2"})
	s.AddClient(&Passenger{Id: "p1"})

	cases := []struct {
		name   string
		filter *ClientFilter
		want   []string
	}{
		{
			name:   "type and id: found",
			filter: &ClientFilter{Id: "d1", Type: string(ClientTypeDriver), Active: -1, Online: -1},
			want:   []string{"d1"},
		},
		{
			name:   "type and id: not found",
			filter: &ClientFilter{Id: "d9", Type: string(ClientTypeDriver), Active: -1, Online: -1},
			want:   []string{},
		},
		{
			name:   "empty filter returns all clients",
			filter: NewEmptyFilter(),
			want:   []string{"d1", "d2", "p1"},
		},
		{
			name:   "only type given",
			filter: &ClientFilter{Type: string(ClientTypeDriver), Active: -1, Online: -1},
			want:   []string{"d1", "d2"},
		},
		{
			name:   "only id given",
			filter: &ClientFilter{Id: "d1", Active: -1, Online: -1},
			want:   []string{"d1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := s.SearchClients(c.filter)

			ids := util.ArrayMap(got, func(c Client) string { return c.getId() })
			sort.Strings(ids)

			test_utils.AssertEquals(c.want, ids, t, false, c.name)
		})
	}
}
