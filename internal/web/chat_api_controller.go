package web

import (
	"log"
	"net/http"
	"strconv"

	"darbelis.eu/taxi/internal/dto"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/repository"
	"darbelis.eu/taxi/internal/state"
	"github.com/gin-gonic/gin"
)

type ChatApiController struct {
	mainState             *state.MainState
	chatMessageRepository repository.ChatMessageRepositoryContract
}

func NewChatApiController(mainState *state.MainState, chatMessageRepository repository.ChatMessageRepositoryContract) *ChatApiController {
	return &ChatApiController{mainState: mainState, chatMessageRepository: chatMessageRepository}
}

// ChatLog handles GET /api/chat/{client_id}/log?with=<otherClientId>
// REST counterpart of handler.ChatMessageHandler.HandleClientRequestsChatLog.
// TODO refactor uri structure
func (c *ChatApiController) ChatLog(ctx *gin.Context) {
	clientId := ctx.Param("client_id")
	otherId := ctx.Query("with")

	if otherId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "with is required"})
		return
	}

	conversation := c.chatMessageRepository.GetConversation(clientId, otherId)

	sender := c.mainState.GetFirstClient(clientId)

	if sender == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "sender not found"})
		return
	}
	receiver := c.mainState.GetFirstClient(otherId)
	if receiver == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
		return
	}

	var chatMessages []messages.ChatMessage
	for _, msg := range conversation {
		chatMessage := transformChatDtoToMessage(msg)
		if chatMessage.SenderId == sender.GetId() {
			chatMessage.SenderName = sender.GetName()
			chatMessage.ReceiverName = receiver.GetName()
		} else {
			chatMessage.SenderName = receiver.GetName()
			chatMessage.ReceiverName = sender.GetName()
		}
		chatMessages = append(chatMessages, chatMessage)
	}

	ctx.JSON(http.StatusOK, messages.ChatLog{
		ParticipantId1: clientId,
		ParticipantId2: otherId,
		Log:            chatMessages,
	})
}

// Talkers handles GET /api/talkers
// REST counterpart of handler.ChatMessageHandler.HandleClientRequestsTalkersInfos.
// Unlike the WS handler, there is no authenticated clientId to exclude from the
// list here, so this returns every online client.
func (c *ChatApiController) Talkers(ctx *gin.Context) {
	filter := state.NewEmptyFilter()
	filter.Online = 1

	clients := c.mainState.SearchClients(filter)
	log.Println("ChatApiController.Talkers clients count: " + strconv.Itoa(len(clients)))

	var talkerInfos []messages.TalkerInfo
	for _, cl := range clients {
		name := cl.GetName()
		if name == "" {
			name = "anonymous"
		}

		talkerInfos = append(talkerInfos, messages.TalkerInfo{
			ClientId:   cl.GetId(),
			ClientType: string(cl.GetType()),
			Name:       name,
		})
	}
	log.Println("ChatApiController.Talkers talkers count: " + strconv.Itoa(len(talkerInfos)))

	ctx.JSON(http.StatusOK, messages.TalkersInfos{List: talkerInfos})
}

func transformChatDtoToMessage(dtoMessage *dto.ChatMessage) messages.ChatMessage {
	return messages.ChatMessage{
		Id:           dtoMessage.Id,
		SenderId:     dtoMessage.SenderId,
		ReceiverId:   dtoMessage.ReceiverId,
		Text:         dtoMessage.Text,
		SentTime:     dtoMessage.SentTime,
		ReceivedTime: dtoMessage.ReceivedTime,
	}
}
