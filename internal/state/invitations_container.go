package state

import "sync"

type InvitationsContainer struct {
	// TODO more efficient data structure after MVP
	// TODO solve storage to database and cleaning memory after MVP
	invitations map[string]*Invitation

	invitationsLock sync.Mutex
}

func NewInvitationsContainer() *InvitationsContainer {
	return &InvitationsContainer{
		invitations:     make(map[string]*Invitation),
		invitationsLock: sync.Mutex{},
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

func (c *InvitationsContainer) Add(invitation Invitation) {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	c.invitations[invitation.Id] = &invitation
}

func (c *InvitationsContainer) Remove(id string) {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	delete(c.invitations, id)
}

func (c *InvitationsContainer) GetInvitationsSnapshot(filter *InvitationsFilter) []*Invitation {
	c.invitationsLock.Lock()
	defer c.invitationsLock.Unlock()

	var invitationsSnapshot []*Invitation

	for _, invitation := range c.invitations {
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
