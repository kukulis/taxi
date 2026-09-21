package state

import (
	"testing"

	"darbelis.eu/taxi/internal/events"
)

func TestMainState_DriverRegistrationRoundTrip(t *testing.T) {
	s := NewMainState()

	registered := events.NewClientRegisteredEvent("driver-1", events.WithRegisteredClientType(events.ClientTypeDriver))
	s.handleDriverRegistered(registered)

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

	unregistered := events.NewClientUnregisteredEvent("driver-1", events.WithUnregisteredClientType(events.ClientTypeDriver))
	s.handleDriverUnregistered(unregistered)

	snapshot = s.GetDriversSnapshot()
	if len(snapshot) != 1 {
		t.Fatalf("GetDriversSnapshot() after unregistration = %d drivers, want 1 (record kept, just inactive)", len(snapshot))
	}
	if snapshot[0].Active {
		t.Errorf("driver Active = true, want false after unregistration")
	}
}
