package handler

import (
	"log/slog"

	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
)

type ChatMessageHandler struct {
	mainState         *state.MainState
	driversHub        ws.HubInterface
	passengersHub     ws.HubInterface
	logger            *slog.Logger
	dispatcher        *util.Dispatcher
	addressWorkerPool *nominatim.AddressWorkerPool
}

func NewChatMessageHandler(
	mainState *state.MainState,
	driversHub ws.HubInterface,
	passengersHub ws.HubInterface,
	dispatcher *util.Dispatcher,
	addressWorkerPool *nominatim.AddressWorkerPool,
) *ChatMessageHandler {
	return &ChatMessageHandler{
		mainState:         mainState,
		driversHub:        driversHub,
		passengersHub:     passengersHub,
		logger:            slog.Default().With("handler", "ChatMessageHandler"),
		dispatcher:        dispatcher,
		addressWorkerPool: addressWorkerPool,
	}
}

func (cmh *ChatMessageHandler) HandleChatMessage(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatMessage)
}

func (cmh *ChatMessageHandler) HandleChatLog(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatLog)
}

func (cmh *ChatMessageHandler) HandleChatInfo(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatInfo)
}

func (cmh *ChatMessageHandler) HandleChatInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatInfos)
}

func (cmh *ChatMessageHandler) HandleClientRequestsChatsInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ClientRequestsChatsInfos)
}

func (cmh *ChatMessageHandler) HandleClientRequestsChatLog(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ClientRequestsChatLog)
}

func (cmh *ChatMessageHandler) HandleTalkerInfo(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.TalkerInfo)
}

func (cmh *ChatMessageHandler) HandleTalkersInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.TalkersInfos)
}

func (cmh *ChatMessageHandler) HandleClientRequestsTalkersInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ClientRequestsTalkersInfos)
}
