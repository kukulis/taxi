package state

import (
	"cmp"
	"slices"
	"testing"

	"darbelis.eu/taxi/pkg/test_utils"
	"darbelis.eu/taxi/pkg/util"
)

func TestInvitationsContainer_AddRepeated(t *testing.T) {
	invitationsContainer := NewInvitationsContainer()

	// try to remove from empty
	invitationsContainer.Remove("0")

	invitationsContainer.Add(
		&Invitation{
			Id:          "1",
			PassengerId: "pass1",
			DriverId:    "driver1",
			Status:      InvitationStatusPending,
		},
	)
	invitationsContainer.Add(
		&Invitation{
			Id:          "2",
			PassengerId: "pass1",
			DriverId:    "driver1",
			Status:      InvitationStatusRejected,
		},
	)

	invitationsContainer.Add(
		&Invitation{
			Id:          "3",
			PassengerId: "pass2",
			DriverId:    "driver1",
			Status:      InvitationStatusAccepted,
		},
	)

	invitationsContainer.Add(
		&Invitation{
			Id:          "4",
			PassengerId: "pass2",
			DriverId:    "driver2",
			Status:      InvitationStatusPending,
		},
	)

	invitationsContainer.Add(
		&Invitation{
			Id:          "5",
			PassengerId: "pass1",
			DriverId:    "driver2",
			Status:      InvitationStatusPending,
		},
	)

	invitations := invitationsContainer.GetInvitationsSnapshot(&InvitationsFilter{
		PassengerId: "pass1",
		DriverId:    "driver1",
	})

	slices.SortFunc(invitations, func(a, b *Invitation) int {
		return cmp.Compare(a.Id, b.Id)
	})

	if len(invitations) != 2 {
		t.Errorf("inviations count = %d, want 2", len(invitations))
	}

	if invitations[0].Id != "1" {
		t.Errorf("inviation[0].Id = %q, want %q", invitations[0].Id, "1")
	}

	if invitations[1].Id != "2" {
		t.Errorf("inviation[1].Id = %q, want %q", invitations[1].Id, "2")
	}

	inviations2 := invitationsContainer.GetInvitationsSnapshot(&InvitationsFilter{
		PassengerId: "pass2",
		DriverId:    "driver1",
	})

	if len(inviations2) != 1 {
		t.Fatalf("Invitations count = %d, expected 1", len(inviations2))
	}

	inviationsDriver1 := invitationsContainer.GetInvitationsSnapshot(&InvitationsFilter{
		DriverId: "driver1",
	})

	if len(inviationsDriver1) != 3 {
		t.Fatalf("Invitations count = %d, expected 3", len(inviationsDriver1))
	}

	slices.SortFunc(inviationsDriver1, InvitationComparatorById)

	idsDriver1 := util.ArrayMap(inviationsDriver1, GetInvitationId)

	test_utils.AssertEquals([]string{"1", "2", "3"}, idsDriver1, t, false, "Invitations for driver 1")

	invitationsAcceptedAndRejected := invitationsContainer.GetInvitationsSnapshot(
		&InvitationsFilter{
			Statuses: []string{InvitationStatusRejected, InvitationStatusAccepted},
		})

	slices.SortFunc(invitationsAcceptedAndRejected, InvitationComparatorById)

	test_utils.AssertEquals([]string{"2", "3"}, util.ArrayMap(invitationsAcceptedAndRejected, GetInvitationId), t, false, "Invitations rejected and accepted")

	invitationsContainer.Remove("3")
	invitationsContainer.Remove("4")

	// after removal

	invitationsPending := invitationsContainer.GetInvitationsSnapshot(
		&InvitationsFilter{
			Statuses: []string{InvitationStatusPending},
		})
	slices.SortFunc(invitationsPending, InvitationComparatorById)
	test_utils.AssertEquals([]string{"1", "5"}, util.ArrayMap(invitationsPending, GetInvitationId), t, false, "Pending invitations after removal")

	diver1AfterDelete := invitationsContainer.GetInvitationsSnapshot(
		&InvitationsFilter{
			DriverId: "driver1",
		})

	slices.SortFunc(diver1AfterDelete, InvitationComparatorById)
	test_utils.AssertEquals([]string{"1", "2"}, util.ArrayMap(diver1AfterDelete, GetInvitationId), t, false, "Driver 1 after deletion")

	pass1AfterDelete := invitationsContainer.GetInvitationsSnapshot(
		&InvitationsFilter{
			PassengerId: "pass1",
		})

	slices.SortFunc(pass1AfterDelete, InvitationComparatorById)
	test_utils.AssertEquals([]string{"1", "2", "5"}, util.ArrayMap(pass1AfterDelete, GetInvitationId), t, false, "Driver 1 after deletion")

	// lets remove last related to driver2
	invitationsContainer.Remove("5")

	// test no panic
	invitationsContainer.Remove("5")

	// try to remove never existed
	invitationsContainer.Remove("6")
}
