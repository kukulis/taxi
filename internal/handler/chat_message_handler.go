package handler

import (
	"log/slog"

	"darbelis.eu/taxi/internal/dto"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/repository"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/message_common"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
)

const MessagesScreenSize = 10

type ChatMessageHandler struct {
	mainState             *state.MainState
	driversHub            ws.HubInterface
	passengersHub         ws.HubInterface
	logger                *slog.Logger
	dispatcher            *util.Dispatcher
	addressWorkerPool     *nominatim.AddressWorkerPool
	chatMessageRepository repository.ChatMessageRepositoryContract
}

func NewChatMessageHandler(
	mainState *state.MainState,
	driversHub ws.HubInterface,
	passengersHub ws.HubInterface,
	dispatcher *util.Dispatcher,
	addressWorkerPool *nominatim.AddressWorkerPool,
	chatMessageRepository repository.ChatMessageRepositoryContract,
) *ChatMessageHandler {
	return &ChatMessageHandler{
		mainState:             mainState,
		driversHub:            driversHub,
		passengersHub:         passengersHub,
		logger:                slog.Default().With("handler", "ChatMessageHandler"),
		dispatcher:            dispatcher,
		addressWorkerPool:     addressWorkerPool,
		chatMessageRepository: chatMessageRepository,
	}
}

func (cmh *ChatMessageHandler) searchClient(clientId string) state.Client {
	senders := cmh.mainState.SearchClients(&state.ClientFilter{Id: clientId})

	// we will decide later what to do if we have two clients with same id : driver and passenger

	if len(senders) > 0 {
		return senders[0]
	}

	return nil
}

func (cmh *ChatMessageHandler) getHubByType(clientType state.ClientType) ws.HubInterface {
	if clientType == state.ClientTypeDriver {
		return cmh.driversHub
	}

	if clientType == state.ClientTypePassenger {
		return cmh.passengersHub
	}

	return nil
}

// HandleChatMessage dedicated to be sent from client to server and from server to client
func (cmh *ChatMessageHandler) HandleChatMessage(clientId string, payload message_common.MessageInterface) {
	chatMessage := payload.(messages.ChatMessage)

	if clientId != chatMessage.SenderId {
		cmh.logger.Warn("The received message sender " + chatMessage.SenderId + " does not match clientId " + clientId)
		return
	}

	sender := cmh.searchClient(chatMessage.SenderId)
	if sender == nil {
		cmh.logger.Warn("Sender not found " + chatMessage.SenderId)
		return
	}

	receiver := cmh.searchClient(chatMessage.ReceiverId)
	if receiver == nil {
		cmh.logger.Warn("Receiver not found " + chatMessage.ReceiverId)
		return
	}

	dtoMessage := transformChatMessageToDto(chatMessage)
	successfulAdd := cmh.chatMessageRepository.Add(dtoMessage)
	if !successfulAdd {
		cmh.logger.Warn("The message with id " + dtoMessage.Id + " already in the repository")
		return
	}

	cmh.chatMessageRepository.CutConversation(dtoMessage.SenderId, dtoMessage.ReceiverId, MessagesScreenSize)

	hub := cmh.getHubByType(receiver.GetType())
	if hub == nil {
		cmh.logger.Warn("Weird, the receiver hub is nil; receiver type: " + string(receiver.GetType()))
		return
	}

	hub.SendMessage(ws.ClientMessage{
		ClientId: chatMessage.ReceiverId,
		Message:  chatMessage,
	})
}

// ============= chat info ===================================================
// not sure if it wouldn't be better to use REST API for these cases

// HandleClientRequestsChatsInfos dedicated to be sent from client to server
func (cmh *ChatMessageHandler) HandleClientRequestsChatsInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ClientRequestsChatsInfos)

	client := cmh.searchClient(clientId)

	if client == nil {
		cmh.logger.Warn("Client not found " + clientId)
		return
	}

	hub := cmh.getHubByType(client.GetType())

	if hub == nil {
		cmh.logger.Error("Could not find hub by a client type " + string(client.GetType()))
		return
	}

	chatPartners := cmh.chatMessageRepository.GetChatPartners(clientId)

	var chatInfos []messages.ChatInfo
	for _, partnerId := range chatPartners {
		chatInfo := messages.ChatInfo{
			ParticipantId1: clientId,
			ParticipantId2: partnerId,
		}

		chatInfos = append(chatInfos, chatInfo)
	}

	hub.SendMessage(ws.ClientMessage{
		ClientId: clientId,
		Message:  messages.ChatInfos{List: chatInfos},
	})

}

// HandleChatInfo dedicated to send from server to client
func (cmh *ChatMessageHandler) HandleChatInfo(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatInfo)
}

// HandleChatInfos dedicated to send from server to client
func (cmh *ChatMessageHandler) HandleChatInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatInfos)
}

// ============= chat log ===================================================

// HandleClientRequestsChatLog dedicated to be sent from client to server
func (cmh *ChatMessageHandler) HandleClientRequestsChatLog(clientId string, payload message_common.MessageInterface) {
	client := cmh.searchClient(clientId)

	if client == nil {
		cmh.logger.Warn("could not find client by id " + clientId)
		return
	}

	request := payload.(messages.ClientRequestsChatLog)
	if request.SenderId != clientId {
		cmh.logger.Warn("unauthorized attempt to get foreign chat log " + clientId + " requested (" + request.SenderId + ", " + request.ReceiverId + ")")
	}

	hub := cmh.getHubByType(client.GetType())

	if hub == nil {
		cmh.logger.Error("Could not find hub by a client type " + string(client.GetType()))
		return
	}

	conversation := cmh.chatMessageRepository.GetConversation(request.SenderId, request.ReceiverId)

	var chatMessages []messages.ChatMessage
	for _, msg := range conversation {
		chatMessages = append(chatMessages, transformChatDtoToMessage(msg))
	}

	hub.SendMessage(ws.ClientMessage{
		ClientId: clientId,
		Message: messages.ChatLog{
			ParticipantId1: request.SenderId,
			ParticipantId2: request.ReceiverId,
			Log:            chatMessages,
		},
	})
}

// HandleChatLog dedicated to send from server to client
func (cmh *ChatMessageHandler) HandleChatLog(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ChatLog)
}

// ============= talker info ===================================================

// HandleClientRequestsTalkersInfos dedicated to be sent from client to server
func (cmh *ChatMessageHandler) HandleClientRequestsTalkersInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.ClientRequestsTalkersInfos)

	client := cmh.searchClient(clientId)

	if client == nil {
		cmh.logger.Warn("could not find client by id " + clientId)
		return
	}

	hub := cmh.getHubByType(client.GetType())

	if hub == nil {
		cmh.logger.Error("Could not find hub by a client type " + string(client.GetType()))
		return
	}

	// just return all clients, transformed in to messages.TalkerInfo
	clients := cmh.mainState.SearchClients(&state.ClientFilter{Online: 1})

	var talkerInfos []messages.TalkerInfo

	for _, cl := range clients {
		if cl.GetId() == client.GetId() {
			continue
		}

		name := cl.GetName()

		if name == "" {
			name = "anonymous"
		}

		talkerInfo := messages.TalkerInfo{
			ClientId:   cl.GetId(),
			ClientType: string(cl.GetType()),
			Name:       name,
		}

		talkerInfos = append(talkerInfos, talkerInfo)
	}

	hub.SendMessage(ws.ClientMessage{
		ClientId: clientId,
		Message: messages.TalkersInfos{
			List: talkerInfos,
		},
	})
}

// HandleTalkerInfo dedicated to send from server to client
func (cmh *ChatMessageHandler) HandleTalkerInfo(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.TalkerInfo)
}

// HandleTalkersInfos dedicated to send from server to client
func (cmh *ChatMessageHandler) HandleTalkersInfos(clientId string, payload message_common.MessageInterface) {
	_ = payload.(messages.TalkersInfos)
}

func transformChatMessageToDto(message messages.ChatMessage) *dto.ChatMessage {
	return &dto.ChatMessage{
		Id:           message.Id,
		SenderId:     message.SenderId,
		ReceiverId:   message.ReceiverId,
		Text:         message.Text,
		SentTime:     message.SentTime,
		ReceivedTime: message.ReceivedTime,
	}
}
func transformChatDtoToMessage(dto *dto.ChatMessage) messages.ChatMessage {
	return messages.ChatMessage{
		Id:           dto.Id,
		SenderId:     dto.SenderId,
		ReceiverId:   dto.ReceiverId,
		Text:         dto.Text,
		SentTime:     dto.SentTime,
		ReceivedTime: dto.ReceivedTime,
	}
}
