package di

import (
	"darbelis.eu/taxi/internal/web"
	"darbelis.eu/taxi/internal/ws"
)

var webControllerInstance *web.WebController = nil
var wsControllerInstance *web.WsController = nil

var driversHubInstance *ws.Hub = nil
var passengersHubInstance *ws.Hub = nil

var driversMessagesWsHandler ws.MessageHandler = nil
var passengersMessagesWsHandler ws.MessageHandler = nil

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
		driversHubInstance = ws.NewHub(GetDriversWsHandler())
	}

	return driversHubInstance
}

func GetPassengersHub() *ws.Hub {
	if passengersHubInstance == nil {
		passengersHubInstance = ws.NewHub(GetPassengersWsHandler())
	}

	return passengersHubInstance
}

func GetDriversWsHandler() ws.MessageHandler {
	if driversMessagesWsHandler == nil {
		driversMessagesWsHandler = &ws.SimpleMessageHandler{}
	}

	return driversMessagesWsHandler
}

func GetPassengersWsHandler() ws.MessageHandler {
	if passengersMessagesWsHandler == nil {
		passengersMessagesWsHandler = &ws.SimpleMessageHandler{}
	}

	return passengersMessagesWsHandler
}
