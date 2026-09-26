package state

import (
	"cmp"
	"log"
	"time"

	"darbelis.eu/taxi/pkg/util"
)

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

	// 9% hack
	Driver    *Driver
	Passenger *Passenger
}

func NewInvitation(id string) *Invitation {
	return &Invitation{Id: id, Status: InvitationStatusPending}
}

func InvitationComparatorById(a *Invitation, b *Invitation) int {
	return cmp.Compare(a.Id, b.Id)
}

func GetInvitationId(i *Invitation) string {
	return i.Id
}

func GetInvitationPassengerId(i *Invitation) string {
	return i.PassengerId
}

func GetInvitationDriverId(i *Invitation) string {
	return i.DriverId
}

func AssignPassengersToInvitations(invitations []*Invitation, passengers []*Passenger) {
	util.ApplyWhenMatch(
		invitations,
		passengers,
		GetInvitationPassengerId,
		GetPassengerId,
		func(i *Invitation, p *Passenger) { i.Passenger = p },
		func(i *Invitation) {
			log.Printf("invitation references unknown passenger %v", i)
		},
	)
}

func AssignDriversToInvitations(invitations []*Invitation, drivers []*Driver) {
	util.ApplyWhenMatch(invitations,
		drivers,
		GetInvitationDriverId,
		GetDriverId,
		func(i *Invitation, d *Driver) { i.Driver = d },
		func(i *Invitation) {
			log.Printf("invitation references unknown driver %v", i)
		},
	)
}

func (i *Invitation) GetDriverLat() float64 {
	if i.Driver == nil {
		return 0
	}

	return i.Driver.Lat
}

func (i *Invitation) GetDriverLon() float64 {
	if i.Driver == nil {
		return 0
	}

	return i.Driver.Lon
}

func (i *Invitation) GetPassengerLat() float64 {
	if i.Passenger == nil {
		return 0
	}

	return i.Passenger.Lat
}

func (i *Invitation) GetPassengerLon() float64 {
	if i.Passenger == nil {
		return 0
	}

	return i.Passenger.Lon
}
