package state

import (
	"sync"
)

type InvitationsContainer struct {
	invitations              map[string]*Invitation
	invitationsByDriverId    map[string][]*Invitation
	invitationsByPassengerId map[string][]*Invitation

	invitationsLock sync.Mutex
}

func NewInvitationsContainer() *InvitationsContainer {
	return &InvitationsContainer{
		invitations:              make(map[string]*Invitation),
		invitationsByDriverId:    make(map[string][]*Invitation),
		invitationsByPassengerId: make(map[string][]*Invitation),
		invitationsLock:          sync.Mutex{},
	}
}

func (c *InvitationsContainer) GetSnapshot(id string) *Invitation {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	invitation, ok := c.invitations[id]

	if !ok {
		return nil
	}

	invitationCopy := *invitation

	return &invitationCopy
}

func (c *InvitationsContainer) Add(invitation *Invitation) {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	c.invitations[invitation.Id] = invitation

	if _, ok := c.invitationsByPassengerId[invitation.PassengerId]; !ok {
		c.invitationsByPassengerId[invitation.PassengerId] = []*Invitation{}
	}

	if _, ok := c.invitationsByPassengerId[invitation.DriverId]; !ok {
		c.invitationsByPassengerId[invitation.DriverId] = []*Invitation{}
	}

	c.invitationsByPassengerId[invitation.PassengerId] = append(c.invitationsByPassengerId[invitation.PassengerId], invitation)
	c.invitationsByDriverId[invitation.DriverId] = append(c.invitationsByDriverId[invitation.DriverId], invitation)
}

func (c *InvitationsContainer) Remove(id string) {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	invitation, ok := c.invitations[id]

	if !ok {
		return
	}

	passengerId := invitation.PassengerId
	driverId := invitation.DriverId

	c.invitationsByPassengerId[passengerId] = DeleteInvitationFromSlice(c.invitationsByPassengerId[passengerId], invitation.Id)
	// check if this was the last invitation for a passenger
	if len(c.invitationsByPassengerId[passengerId]) == 0 {
		delete(c.invitationsByPassengerId, passengerId)
	}

	c.invitationsByDriverId[driverId] = DeleteInvitationFromSlice(c.invitationsByDriverId[driverId], invitation.Id)
	// check if this was the last invitation for a driver
	if len(c.invitationsByDriverId[driverId]) == 0 {
		delete(c.invitationsByDriverId, driverId)
	}

	delete(c.invitations, id)

}

func (c *InvitationsContainer) GetInvitationsSnapshot(filter *InvitationsFilter) []*Invitation {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	if filter.PassengerId != "" {
		return c.getInvitationsSnapshotByPassenger(filter)
	}

	if filter.DriverId != "" {
		return c.getInvitationsSnapshotByDriver(filter)
	}

	var invitationsSnapshot []*Invitation
	for _, invitation := range c.invitations {
		if filter.Match(*invitation) {
			invitationCopy := *invitation
			invitationsSnapshot = append(invitationsSnapshot, &invitationCopy)
		}
	}

	return invitationsSnapshot
}

func (c *InvitationsContainer) getInvitationsSnapshotByPassenger(filter *InvitationsFilter) []*Invitation {
	var invitationsSnapshot []*Invitation
	for _, invitation := range c.invitationsByPassengerId[filter.PassengerId] {
		if filter.Match(*invitation) {
			invitationCopy := *invitation
			invitationsSnapshot = append(invitationsSnapshot, &invitationCopy)
		}
	}

	return invitationsSnapshot
}

func (c *InvitationsContainer) getInvitationsSnapshotByDriver(filter *InvitationsFilter) []*Invitation {
	var invitationsSnapshot []*Invitation
	for _, invitation := range c.invitationsByDriverId[filter.DriverId] {
		if filter.Match(*invitation) {
			invitationCopy := *invitation
			invitationsSnapshot = append(invitationsSnapshot, &invitationCopy)
		}
	}

	return invitationsSnapshot
}

func (c *InvitationsContainer) Update(id string, mutator func(invitation *Invitation)) {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	mutator(c.invitations[id])
}

func DeleteInvitationFromSlice(invitations []*Invitation, invitationId string) []*Invitation {

	foundIndex := -1
	for i, invitation := range invitations {
		if invitationId == invitation.Id {
			foundIndex = i
			break
		}
	}

	if foundIndex < 0 {
		return invitations
	}

	// invitation found

	// put the last element instead of deleted if it is not the last element
	if len(invitations)-1 > foundIndex {
		invitations[foundIndex] = invitations[len(invitations)-1]
	}

	// return shorter slice by one element
	return invitations[0 : len(invitations)-1]
}
