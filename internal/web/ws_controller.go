package web

import (
	"log"

	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
	"github.com/gin-gonic/gin"
)

type WsController struct {
	driversHub    *ws.Hub
	passengersHub *ws.Hub
}

func NewWsController(driversHub *ws.Hub, passengersHub *ws.Hub) *WsController {
	return &WsController{driversHub: driversHub, passengersHub: passengersHub}
}

func (c *WsController) ServeDriversWebSocketConnection(ctx *gin.Context) {
	clientId, err := util.RandomHash(8)
	if err != nil {
		log.Println("Failed to create a client id", err)
		ctx.AbortWithError(500, err)
		return
	}
	c.driversHub.RegisterWebSocketClient(clientId, ctx.Writer, ctx.Request)
}

func (c *WsController) ServePassengersWebSocketConnection(ctx *gin.Context) {
	clientId, err := util.RandomHash(8)
	if err != nil {
		log.Println("Failed to create a client id", err)
		ctx.AbortWithError(500, err)
		return
	}
	c.passengersHub.RegisterWebSocketClient(clientId, ctx.Writer, ctx.Request)
}
