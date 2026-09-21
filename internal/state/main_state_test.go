package state

import (
	"testing"
	"time"

	"darbelis.eu/taxi/internal/events"
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

func TestMainState_DriverRegistrationRoundTrip(t *testing.T) {
	s := NewMainState(util.RealClock{})
	go s.HandleDedicatedEvents()

	s.AddDedicatedEvent(events.NewClientRegisteredEvent("driver-1", events.WithRegisteredClientType(events.ClientTypeDriver)))
	waitFor(t, 200*time.Millisecond, func() bool {
		driver := s.GetDriverById("driver-1")
		return driver != nil && driver.Active
	}, "driver-1 was never registered as active")

	snapshot := s.GetDriversSnapshot()
	if len(snapshot) != 1 {
		t.Fatalf("GetDriversSnapshot() after registration = %d drivers, want 1", len(snapshot))
	}
	if snapshot[0].Id != "driver-1" {
		t.Errorf("driver Id = %q, want %q", snapshot[0].Id, "driver-1")
	}
	if !snapshot[0].Active {
		t.Errorf("driver Active = false, want true after registration")
	}

	s.AddDedicatedEvent(events.NewClientUnregisteredEvent("driver-1", events.WithUnregisteredClientType(events.ClientTypeDriver)))
	waitFor(t, 200*time.Millisecond, func() bool {
		driver := s.GetDriverById("driver-1")
		return driver != nil && !driver.Active
	}, "driver-1 was never marked inactive")

	snapshot = s.GetDriversSnapshot()
	if len(snapshot) != 1 {
		t.Fatalf("GetDriversSnapshot() after unregistration = %d drivers, want 1 (record kept, just inactive)", len(snapshot))
	}
	if snapshot[0].Active {
		t.Errorf("driver Active = true, want false after unregistration")
	}
}

func TestMainState_PassengerRegistrationRoundTrip(t *testing.T) {
	s := NewMainState(util.RealClock{})
	go s.HandleDedicatedEvents()

	s.AddDedicatedEvent(events.NewClientRegisteredEvent("passenger-1", events.WithRegisteredClientType(events.ClientTypePassenger)))
	waitFor(t, 200*time.Millisecond, func() bool {
		passenger := findPassenger(s, "passenger-1")
		return passenger != nil && passenger.Active
	}, "passenger-1 was never registered as active")

	snapshot := s.GetPassengersSnapshot()
	if len(snapshot) != 1 {
		t.Fatalf("GetPassengersSnapshot() after registration = %d passengers, want 1", len(snapshot))
	}
	if snapshot[0].Id != "passenger-1" {
		t.Errorf("passenger Id = %q, want %q", snapshot[0].Id, "passenger-1")
	}
	if !snapshot[0].Active {
		t.Errorf("passenger Active = false, want true after registration")
	}

	s.AddDedicatedEvent(events.NewClientUnregisteredEvent("passenger-1", events.WithUnregisteredClientType(events.ClientTypePassenger)))
	waitFor(t, 200*time.Millisecond, func() bool {
		passenger := findPassenger(s, "passenger-1")
		return passenger != nil && !passenger.Active
	}, "passenger-1 was never marked inactive")

	snapshot = s.GetPassengersSnapshot()
	if len(snapshot) != 1 {
		t.Fatalf("GetPassengersSnapshot() after unregistration = %d passengers, want 1 (record kept, just inactive)", len(snapshot))
	}
	if snapshot[0].Active {
		t.Errorf("passenger Active = true, want false after unregistration")
	}
}
