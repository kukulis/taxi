package messages

import (
	"time"

	"darbelis.eu/taxi/pkg/message_common"
)

type ChatMessage struct {
	Id           string    `json:"id"`
	SenderId     string    `json:"sender_id"`
	ReceiverId   string    `json:"receiver_id"`
	Text         string    `json:"text"`
	SentTime     time.Time `json:"sent_time"`
	ReceivedTime time.Time `json:"received_time"`
	SenderName   string    `json:"sender_name"`
	ReceiverName string    `json:"receiver_name"`
}

func (ChatMessage) GetMessageType() message_common.MessageType {
	return MessageTypeChatMessage
}

type ChatLog struct {
	Id             string        `json:"id"`
	ParticipantId1 string        `json:"participant_id_1"`
	ParticipantId2 string        `json:"participant_id_2"`
	Log            []ChatMessage `json:"log"`
	StartedAt      time.Time     `json:"started_at"`
}

func (ChatLog) GetMessageType() message_common.MessageType {
	return MessageTypeChatLog
}

type ChatInfo struct {
	Id             string    `json:"id"`
	ParticipantId1 string    `json:"participant_id_1"`
	ParticipantId2 string    `json:"participant_id_2"`
	StartedAt      time.Time `json:"started_at"`
}

func (ChatInfo) GetMessageType() message_common.MessageType {
	return MessageTypeChatInfo
}

type ChatInfos struct {
	List []ChatInfo `json:"list"`
}

func (ChatInfos) GetMessageType() message_common.MessageType {
	return MessageTypeChatInfos
}

type ClientRequestsChatsInfos struct {
}

func (ClientRequestsChatsInfos) GetMessageType() message_common.MessageType {
	return MessageTypeClientRequestsChatsInfos
}

type ClientRequestsChatLog struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
}

func (ClientRequestsChatLog) GetMessageType() message_common.MessageType {
	return MessageTypeClientRequestsChatLog
}

type TalkerInfo struct {
	ClientId   string `json:"client_id"`
	ClientType string `json:"client_type"`
	Name       string `json:"name"`
}

func (TalkerInfo) GetMessageType() message_common.MessageType {
	return MessageTypeTalkerInfo
}

type TalkersInfos struct {
	List []TalkerInfo `json:"list"`
}

func (TalkersInfos) GetMessageType() message_common.MessageType {
	return MessageTypeTalkersInfos
}

type ClientRequestsTalkersInfos struct {
}

func (ClientRequestsTalkersInfos) GetMessageType() message_common.MessageType {
	return MessageTypeClientRequestsTalkersInfos
}
