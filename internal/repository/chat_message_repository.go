package repository

import (
	"sync"

	"darbelis.eu/taxi/internal/dto"
)

type ChatMessageRepositoryContract interface {
	Add(chatMessage *dto.ChatMessage)
	Remove(messageId string) bool
	Get(messageId string) *dto.ChatMessage
	GetConversation(participant1, participant2 string) []*dto.ChatMessage
	GetChatPartners(participantId string) []string
	CutConversation(participant1, participant2 string, size int)
}

type ChatMessageRepositoryMemory struct {
	chatMessagesById map[string]*dto.ChatMessage
	chatPartners     map[string][]string
	conversations    map[string][]*dto.ChatMessage
	messagesLock     sync.Mutex
}

func NewChatMessageRepositoryMemory() *ChatMessageRepositoryMemory {
	return &ChatMessageRepositoryMemory{
		chatMessagesById: make(map[string]*dto.ChatMessage),
		conversations:    make(map[string][]*dto.ChatMessage),
		chatPartners:     make(map[string][]string),
		messagesLock:     sync.Mutex{},
	}
}

func (r *ChatMessageRepositoryMemory) Add(chatMessage *dto.ChatMessage) {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	r.chatMessagesById[chatMessage.Id] = chatMessage

	pairKey := chatMessage.GetPairKey()
	r.conversations[pairKey] = append(r.conversations[pairKey], chatMessage)

	// The message is first in the conversation
	if len(r.conversations[pairKey]) == 1 {
		r.chatPartners[chatMessage.SenderId] = append(r.chatPartners[chatMessage.SenderId], chatMessage.ReceiverId)
		r.chatPartners[chatMessage.ReceiverId] = append(r.chatPartners[chatMessage.ReceiverId], chatMessage.SenderId)
	}
}

func (r *ChatMessageRepositoryMemory) Remove(messageId string) bool {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	message, ok := r.chatMessagesById[messageId]
	if !ok {
		return false
	}

	delete(r.chatMessagesById, messageId)

	pairKey := message.GetPairKey()

	conversation := r.conversations[pairKey]
	foundIndex := -1
	for i := 0; i < len(conversation); i++ {
		if conversation[i].Id == message.Id {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		conversation = append(conversation[:foundIndex], conversation[foundIndex+1:]...)
		r.conversations[pairKey] = conversation

		if len(conversation) == 0 {
			r.removeChatPartner(message.SenderId, message.ReceiverId)
			r.removeChatPartner(message.ReceiverId, message.SenderId)
		}
	}

	return true
}

func (r *ChatMessageRepositoryMemory) removeChatPartner(owner, partner string) {
	partners := r.chatPartners[owner]

	for i, existing := range partners {
		if existing == partner {
			r.chatPartners[owner] = append(partners[:i], partners[i+1:]...)
			return
		}
	}
}

func (r *ChatMessageRepositoryMemory) Get(messageId string) *dto.ChatMessage {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	return r.chatMessagesById[messageId]
}

func (r *ChatMessageRepositoryMemory) GetConversation(participant1, participant2 string) []*dto.ChatMessage {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	pairKey := dto.GetChatMessagePairKey(participant1, participant2)
	snapShot := append([]*dto.ChatMessage{}, r.conversations[pairKey]...)

	return snapShot
}

func (r *ChatMessageRepositoryMemory) GetChatPartners(participantId string) []string {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	snapShot := append([]string{}, r.chatPartners[participantId]...)

	return snapShot
}

// CutConversation leaves 'size' amount of last messages in the conversation
func (r *ChatMessageRepositoryMemory) CutConversation(participant1, participant2 string, size int) {
	r.messagesLock.Lock()
	defer r.messagesLock.Unlock()

	pairKey := dto.GetChatMessagePairKey(participant1, participant2)

	conversation := r.conversations[pairKey]
	l := len(conversation)
	if l > size {
		r.conversations[pairKey] = conversation[l-size : size]
	}
}
