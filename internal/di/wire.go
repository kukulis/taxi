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
