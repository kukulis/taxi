package web

import (
	"log"

	"darbelis.eu/taxi/pkg/ws"
	"github.com/gin-gonic/gin"
)

const (
	ClientCookieName = "gt_taxi_id"
)

type WsController struct {
	driversHub    *ws.Hub
	passengersHub *ws.Hub
}

func NewWsController(driversHub *ws.Hub, passengersHub *ws.Hub) *WsController {
	return &WsController{driversHub: driversHub, passengersHub: passengersHub}
}

func (c *WsController) ServeDriversWebSocketConnection(ctx *gin.Context) {
	clientId, err := RollCookie(ctx)
	if err != nil {
		log.Println("Failed to create a client id", err)
		ctx.AbortWithError(500, err)
		return
	}

	c.driversHub.RegisterWebSocketClient(clientId, ctx.Writer, ctx.Request)
}

func (c *WsController) ServePassengersWebSocketConnection(ctx *gin.Context) {
	clientId, err := RollCookie(ctx)
	if err != nil {
		log.Println("Failed to create a client id", err)
		ctx.AbortWithError(500, err)
		return
	}

	c.passengersHub.RegisterWebSocketClient(clientId, ctx.Writer, ctx.Request)
}
