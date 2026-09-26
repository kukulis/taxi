package handler

import (
	"testing"

	"darbelis.eu/taxi/internal/messages"
	ws2 "darbelis.eu/taxi/pkg/ws"
)

func TestHubMock_RespondsToCoordinatesRequestWithClientCoordinates(t *testing.T) {
	mock := ws2.NewHubMock()
	mock.SendMessageFunc = func(sent ws2.ClientMessage) {
		if _, ok := sent.Message.(messages.ServerRequestsCoordinates); !ok {
			return
		}
		mock.FeedIncomingMessage(ws2.ClientMessage{
			ClientId: sent.ClientId,
			Message:  messages.ClientRespondsCoordinates{Lat: 54.68, Lon: 25.28},
		})
	}

	mock.SendMessage(ws2.ClientMessage{
		ClientId: "driver-1",
		Message:  messages.ServerRequestsCoordinates{},
	})

	got := <-mock.GetIncomingMessagesChannel()

	if got.ClientId != "driver-1" {
		t.Errorf("ClientId = %q, want %q", got.ClientId, "driver-1")
	}

	coords, ok := got.Message.(messages.ClientRespondsCoordinates)
	if !ok {
		t.Fatalf("Message type = %T, want messages.ClientRespondsCoordinates", got.Message)
	}
	if coords.Lat != 54.68 || coords.Lon != 25.28 {
		t.Errorf("coords = %+v, want {Lat:54.68 Lon:25.28}", coords)
	}
}
