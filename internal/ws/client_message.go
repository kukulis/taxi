package ws

import "darbelis.eu/taxi/internal/messages"

type ClientMessage struct {
	ClientId  string
	MessageId string
	Message   messages.MessageInterface
}
