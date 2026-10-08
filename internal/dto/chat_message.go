package dto

import "time"

type ChatMessage struct {
	Id           string
	SenderId     string
	ReceiverId   string
	Text         string
	SentTime     time.Time
	ReceivedTime time.Time
}

func NewChatMessage(id, senderId, receiverId, text string) *ChatMessage {
	return &ChatMessage{
		Id:         id,
		SenderId:   senderId,
		ReceiverId: receiverId,
		Text:       text,
	}
}

func (m *ChatMessage) GetPairKey() string {
	return GetChatMessagePairKey(m.SenderId, m.ReceiverId)
}

func GetChatMessagePairKey(p1, p2 string) string {
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	return p1 + "__" + p2
}
