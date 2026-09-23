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
var driverApiControllerInstance *web.DriverApiController = nil

var driversHubInstance *ws.Hub = nil
var passengersHubInstance *ws.Hub = nil

var driversMessagesWsHandler handler.MessageHandler = nil
var passengersMessagesWsHandler handler.MessageHandler = nil

var mainState *state.MainState = nil

var dispatcher *util.Dispatcher = nil

var taxiDispatcher *handler.TaxiDispatcher = nil

func GetWsController() *web.WsController {
	if wsControllerInstance == nil {
		wsControllerInstance = web.NewWsController(GetDriversHub(), GetPassengersHub())
	}

	return wsControllerInstance
}

func GetWebController() *web.WebController {
	if webControllerInstance == nil {
		webControllerInstance = web.NewWebController(GetMainState())
	}

	return webControllerInstance
}

func GetDriverApiController() *web.DriverApiController {
	if driverApiControllerInstance == nil {
		driverApiControllerInstance = web.NewDriverApiController(GetMainState())
	}

	return driverApiControllerInstance
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
		driversMessagesWsHandler = handler.NewDriversMessageHandler(GetMainState(), GetDriversHub(), GetPassengersHub(), GetDispatcher())
	}

	return driversMessagesWsHandler
}

func GetPassengersWsHandler() handler.MessageHandler {
	if passengersMessagesWsHandler == nil {
		passengersMessagesWsHandler = handler.NewPassengersMessageHandler(GetMainState(), GetDriversHub(), GetPassengersHub())
	}

	return passengersMessagesWsHandler
}

func GetMainState() *state.MainState {
	if mainState == nil {
		mainState = state.NewMainState(util.RealClock{})
	}

	return mainState
}

func GetDispatcher() *util.Dispatcher {
	if dispatcher == nil {
		dispatcher = util.NewDispatcher()
	}

	return dispatcher
}

func GetTaxiDispatcher() *handler.TaxiDispatcher {
	if taxiDispatcher == nil {
		taxiDispatcher = handler.NewTaxiDispatcher(GetMainState(), GetDriversHub(), GetPassengersHub(), util.RealClock{})
	}

	return taxiDispatcher
}
