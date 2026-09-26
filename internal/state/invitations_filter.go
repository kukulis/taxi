package state

type InvitationsFilter struct {
	PassengerId string
	DriverId    string
	Status      string

	Statuses []string
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

	if f.Statuses != nil && len(f.Statuses) != 0 {

		for _, s := range f.Statuses {
			if s == invitation.Status {
				return true
			}
		}
		return false
	}

	return true
}
