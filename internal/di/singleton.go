package di

import (
	"darbelis.eu/taxi/internal/handler"
	"darbelis.eu/taxi/internal/web"
	"darbelis.eu/taxi/internal/ws"
)

var webControllerInstance *web.WebController = nil
var wsControllerInstance *web.WsController = nil

var driversHubInstance *ws.Hub = nil
var passengersHubInstance *ws.Hub = nil

var driversMessagesWsHandler handler.MessageHandler = nil
var passengersMessagesWsHandler handler.MessageHandler = nil

func GetWsController() *web.WsController {
	if wsControllerInstance == nil {
		wsControllerInstance = web.NewWsController(GetDriversHub(), GetPassengersHub())
	}

	return wsControllerInstance
}

func GetWebController() *web.WebController {
	if webControllerInstance == nil {
		webControllerInstance = web.NewWebController()
	}

	return webControllerInstance
}

func GetDriversHub() *ws.Hub {
	if driversHubInstance == nil {
		driversHubInstance = ws.NewHub()
	}

	return driversHubInstance
}

func GetPassengersHub() *ws.Hub {
	if passengersHubInstance == nil {
		passengersHubInstance = ws.NewHub()
	}

	return passengersHubInstance
}

func GetDriversWsHandler() handler.MessageHandler {
	if driversMessagesWsHandler == nil {
		driversMessagesWsHandler = &handler.SimpleMessageHandler{}
	}

	return driversMessagesWsHandler
}

func GetPassengersWsHandler() handler.MessageHandler {
	if passengersMessagesWsHandler == nil {
		passengersMessagesWsHandler = &handler.SimpleMessageHandler{}
	}

	return passengersMessagesWsHandler
}
