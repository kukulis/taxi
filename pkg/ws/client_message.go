package ws

import (
	"darbelis.eu/taxi/pkg/message_common"
)

type ClientMessage struct {
	ClientId  string
	MessageId string
	Message   message_common.MessageInterface
}
