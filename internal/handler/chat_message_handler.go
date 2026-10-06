package handler

import (
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/pkg/message_common"
)

type ChatMessageHandler struct {
	// TODO
}

func NewChatMessageHandler() *ChatMessageHandler {
	return &ChatMessageHandler{}
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
