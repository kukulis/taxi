package state

import (
	"testing"
	"time"

	"darbelis.eu/taxi/pkg/util"
)

// waitFor polls condition until it returns true. If timeout elapses first, it fails the
// test with msg, so the failure is more informative than a generic timeout notice.
func waitFor(t *testing.T, timeout time.Duration, condition func() bool, msg string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatal(msg)
}

func findPassenger(s *MainState, id string) *Passenger {
	for _, passenger := range s.GetPassengersSnapshot() {
		if passenger.Id == id {
			return passenger
		}
	}
	return nil
}

func TestMainState_GetNearestDrivers(t *testing.T) {
	s := NewMainState(util.NewFixedClock(time.Now()))

	const queryLat, queryLon = 54.6872, 25.2797 // Vilnius

	setDriverCoords := func(id string, lat, lon float64) {
		s.CreateDriver(id)
		s.UpdateDriver(id, func(driver *Driver) {
			driver.Lat = lat
			driver.Lon = lon
			driver.CoordinatesReceivedAt = s.Clock.Now()
		})
	}

	setDriverCoords("far", 54.8985, 23.9036)    // Kaunas, ~100km away
	setDriverCoords("near", 54.6892, 25.2799)   // a couple hundred meters away
	setDriverCoords("middle", 54.7500, 25.0000) // somewhere in between

	// never responded with coordinates: must be excluded
	s.CreateDriver("no-coords")

	// inactive with coordinates: must be excluded
	setDriverCoords("inactive", 54.6892, 25.2799)
	s.RemoveDriver("inactive")

	results := s.GetNearestDrivers(queryLat, queryLon)

	if len(results) != 3 {
		t.Fatalf("GetNearestDrivers() returned %d drivers, want 3 (excluding no-coords and inactive): %+v", len(results), results)
	}

	gotOrder := []string{results[0].Driver.Id, results[1].Driver.Id, results[2].Driver.Id}
	wantOrder := []string{"near", "middle", "far"}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Errorf("results[%d].Driver.Id = %q, want %q (order = %v, want %v)", i, gotOrder[i], wantOrder[i], gotOrder, wantOrder)
			break
		}
	}

	for i := 1; i < len(results); i++ {
		if results[i].Distance < results[i-1].Distance {
			t.Errorf("results not sorted ascending by distance: %+v", results)
			break
		}
	}
}
