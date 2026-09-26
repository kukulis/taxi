package handler

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
)

type DriverStatusChangeForClientListener struct {
	driversHub *ws.Hub
}

func NewDriverStatusChangeForClientListener(
	driversHub *ws.Hub,
) *DriverStatusChangeForClientListener {
	return &DriverStatusChangeForClientListener{
		driversHub: driversHub,
	}

}

func (l *DriverStatusChangeForClientListener) Handle(e util.Event) {
	statusChangeEvent := e.(events.DriverStatusChangedEvent)

	l.driversHub.SendMessage(ws.ClientMessage{
		ClientId: statusChangeEvent.DriverId,
		Message:  messages.ServerRequestsCoordinates{},
	})
}
