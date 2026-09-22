package state

type InvitationsFilter struct {
	PassengerId string
	DriverId    string
	Status      string
}

func (f *InvitationsFilter) Match(invitation Invitation) bool {
	if f.DriverId != "" && f.DriverId != invitation.DriverId {
		return false
	}

	if f.PassengerId != "" && f.PassengerId != invitation.PassengerId {
		return false
	}

	if f.Status != "" && f.Status != invitation.Status {
		return false
	}

	return true
}
