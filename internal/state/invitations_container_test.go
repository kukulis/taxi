package state

import "testing"

func TestInvitationsContainer_Add(t *testing.T) {
	invitationsContainer := NewInvitationsContainer()

	invitationsContainer.Add(Invitation{
		Id:          "1",
		PassengerId: "pass1",
		DriverId:    "driver1",
	})

	s := invitationsContainer.GetSnapshot("1")

	if s == nil {
		t.Fatalf("s = nil")
	}

	expectedId := "1"

	if s.Id != expectedId {
		t.Errorf("s.Id = %q, want %q", s.Id, expectedId)
	}
}

func TestInvitationsContainer_AddRepeated(t *testing.T) {
	invitationsContainer := NewInvitationsContainer()

	invitationsContainer.Add(
		Invitation{
			Id:          "1",
			PassengerId: "pass1",
			DriverId:    "driver1",
		},
	)
	invitationsContainer.Add(
		Invitation{
			Id:          "2",
			PassengerId: "pass1",
			DriverId:    "driver1",
		},
	)

	inviations := invitationsContainer.GetInvitationsSnapshot(&InvitationsFilter{
		PassengerId: "pass1",
		DriverId:    "driver1",
	})

	if len(inviations) != 2 {
		t.Errorf("inviations count = %d, want 2", len(inviations))
	}

	if inviations[0].Id != "1" {
		t.Errorf("inviation[0].Id = %q, want %q", inviations[0].Id, "1")
	}

	if inviations[1].Id != "2" {
		t.Errorf("inviation[1].Id = %q, want %q", inviations[1].Id, "2")
	}
}
