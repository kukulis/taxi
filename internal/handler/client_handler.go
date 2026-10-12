package handler

import (
	"log/slog"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/repository"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
)

type ClientHandler struct {
	mainState             *state.MainState
	driversHub            ws.HubInterface
	passengersHub         ws.HubInterface
	logger                *slog.Logger
	dispatcher            *util.Dispatcher
	addressWorkerPool     *nominatim.AddressWorkerPool
	chatMessageRepository repository.ChatMessageRepositoryContract
}

func NewClientHandler(
	mainState *state.MainState,
	driversHub ws.HubInterface,
	passengersHub ws.HubInterface,
	dispatcher *util.Dispatcher,
	addressWorkerPool *nominatim.AddressWorkerPool,
	chatMessageRepository repository.ChatMessageRepositoryContract,
) *ClientHandler {
	return &ClientHandler{
		mainState:             mainState,
		driversHub:            driversHub,
		passengersHub:         passengersHub,
		logger:                slog.Default().With("handler", "ClientHandler"),
		dispatcher:            dispatcher,
		addressWorkerPool:     addressWorkerPool,
		chatMessageRepository: chatMessageRepository,
	}
}

func (ch *ClientHandler) searchClient(clientId string) state.Client {
	filter := state.NewEmptyFilter()
	filter.Id = clientId
	senders := ch.mainState.SearchClients(filter)

	// we will decide later what to do if we have two clients with same id : driver and passenger

	if len(senders) > 0 {
		return senders[0]
	}

	return nil
}

// HandleNameMessage dedicated to be sent from client to server and from server to client
func (ch *ClientHandler) HandleNameMessage(clientId string, payload message_common.MessageInterface) {
	clientName := payload.(*messages.ClientName)

	sender := ch.searchClient(clientId)
	if sender == nil {
		ch.logger.Warn("Sender not found " + clientId)
		return
	}

	ch.mainState.UpdateClient(sender, func(c state.Client) {
		c.SetName(clientName.Name)
	})

}
