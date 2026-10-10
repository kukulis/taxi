package web

import (
	"net/http"

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
func (c *ChatApiController) ChatLog(ctx *gin.Context) {
	clientId := ctx.Param("client_id")
	otherId := ctx.Query("with")

	if otherId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "with is required"})
		return
	}

	conversation := c.chatMessageRepository.GetConversation(clientId, otherId)

	var chatMessages []messages.ChatMessage
	for _, msg := range conversation {
		chatMessages = append(chatMessages, transformChatDtoToMessage(msg))
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
	clients := c.mainState.SearchClients(&state.ClientFilter{Online: 1})

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
