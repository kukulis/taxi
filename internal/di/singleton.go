package di

import (
	"darbelis.eu/taxi/internal/events"
	"darbelis.eu/taxi/internal/handler"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/web"
	"darbelis.eu/taxi/internal/ws"
	"darbelis.eu/taxi/pkg/util"
)

var webControllerInstance *web.WebController = nil
var wsControllerInstance *web.WsController = nil

var driversHubInstance *ws.Hub = nil
var passengersHubInstance *ws.Hub = nil

var driversMessagesWsHandler handler.MessageHandler = nil
var passengersMessagesWsHandler handler.MessageHandler = nil

var mainState *state.MainState = nil

var dispatcher *util.Dispatcher = nil

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
		driversHubInstance = ws.NewHub(GetDispatcher(), events.ClientTypeDriver)
	}

	return driversHubInstance
}

func GetPassengersHub() *ws.Hub {
	if passengersHubInstance == nil {
		passengersHubInstance = ws.NewHub(GetDispatcher(), events.ClientTypePassenger)
	}

	return passengersHubInstance
}

func GetDriversWsHandler() handler.MessageHandler {
	if driversMessagesWsHandler == nil {
		driversMessagesWsHandler = handler.NewDriversMessageHandler(GetMainState())
	}

	return driversMessagesWsHandler
}

func GetPassengersWsHandler() handler.MessageHandler {
	if passengersMessagesWsHandler == nil {
		passengersMessagesWsHandler = handler.NewPassengersMessageHandler(GetMainState())
	}

	return passengersMessagesWsHandler
}

func GetMainState() *state.MainState {
	if mainState == nil {
		mainState = state.NewMainState()
	}

	return mainState
}

func GetDispatcher() *util.Dispatcher {
	if dispatcher == nil {
		dispatcher = util.NewDispatcher()
	}

	return dispatcher
}
