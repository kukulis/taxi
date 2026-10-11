package di

import (
	"darbelis.eu/taxi/internal/handler"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/pkg/message_common"
)

func WireDelegatingDriverMessageHandler(
	delegating *handler.DelegatingMessageHandler,
	d *handler.DriversMessageHandler) *handler.DelegatingMessageHandler {

	delegating.SetHandler(messages.MessageTypeClientRespondsCoordinates, d.HandleClientRespondsCoordinates)
	delegating.SetHandler(messages.MessageTypeClientRespondsInfo, d.HandleClientRespondsInfo)
	delegating.SetHandler(messages.MessageTypeClientResponseError, d.HandleClientResponseError)
	delegating.SetHandler(messages.MessageTypeClientError, d.HandleClientError)

	delegating.SetHandler(messages.MessageTypeDriverCancelsOffer, d.HandleDriverCancelsOffer)
	delegating.SetHandler(messages.MessageTypeDriverAcceptsOffer, d.HandleDriverAcceptsOffer)
	delegating.SetHandler(messages.MessageTypeDriverRejectsOffer, d.HandleDriverRejectsOffer)
	delegating.SetHandler(messages.MessageTypeDriverChangesStatus, d.HandleDriverChangesStatus)
	delegating.SetHandler(message_common.VirtualMessageTypeClientRegistered, d.HandleClientRegistered)
	delegating.SetHandler(message_common.VirtualMessageTypeClientUnregistered, d.HandleClientUnregistered)

	return delegating
}

func WireDelegatingPassengerMessageHandler(
	delegating *handler.DelegatingMessageHandler,
	p *handler.PassengersMessageHandler) *handler.DelegatingMessageHandler {

	delegating.SetHandler(messages.MessageTypeClientRespondsCoordinates, p.HandleClientRespondsCoordinates)
	delegating.SetHandler(messages.MessageTypeClientRespondsInfo, p.HandleClientRespondsInfo)
	delegating.SetHandler(messages.MessageTypeClientResponseError, p.HandleClientResponseError)
	delegating.SetHandler(messages.MessageTypeClientError, p.HandleClientError)

	delegating.SetHandler(messages.MessageTypePassengerInvitesDriver, p.HandlePassengerInvitesDriver)
	delegating.SetHandler(messages.MessageTypePassengerCancelsInvite, p.HandlePassengerCancelsInvite)
	delegating.SetHandler(messages.MessageTypePassengerRequestDriverCoords, p.HandlePassengerRequestDriverCoords)
	delegating.SetHandler(message_common.VirtualMessageTypeClientRegistered, p.HandleClientRegistered)
	delegating.SetHandler(message_common.VirtualMessageTypeClientUnregistered, p.HandleClientUnregistered)

	return delegating
}

func WireChatMessagesHandles(
	delegating *handler.DelegatingMessageHandler,
	c *handler.ChatMessageHandler) *handler.DelegatingMessageHandler {

	delegating.SetHandler(messages.MessageTypeChatMessage, c.HandleChatMessage)
	delegating.SetHandler(messages.MessageTypeChatLog, c.HandleChatLog)
	delegating.SetHandler(messages.MessageTypeChatInfo, c.HandleChatInfo)
	delegating.SetHandler(messages.MessageTypeChatInfos, c.HandleChatInfos)
	delegating.SetHandler(messages.MessageTypeClientRequestsChatsInfos, c.HandleClientRequestsChatsInfos)
	delegating.SetHandler(messages.MessageTypeClientRequestsChatLog, c.HandleClientRequestsChatLog)
	delegating.SetHandler(messages.MessageTypeTalkerInfo, c.HandleTalkerInfo)
	delegating.SetHandler(messages.MessageTypeTalkersInfos, c.HandleTalkersInfos)
	delegating.SetHandler(messages.MessageTypeClientRequestsTalkersInfos, c.HandleClientRequestsTalkersInfos)

	return delegating
}
