package state

import "time"

const (
	InvitationStatusPending   = "pending"
	InvitationStatusAccepted  = "accepted"
	InvitationStatusRejected  = "rejected"
	InvitationStatusCanceled  = "cancelled"
	InvitationStatusDriving   = "driving"
	InvitationStatusCompleted = "completed"
)

type Invitation struct {
	Id          string
	PassengerId string
	DriverId    string
	Status      string
	CreatedAt   time.Time
	AcceptedAt  time.Time
	RejectedAt  time.Time
	CancelledAt time.Time
}

func NewInvitation(id string) *Invitation {
	return &Invitation{Id: id, Status: InvitationStatusPending}
}
