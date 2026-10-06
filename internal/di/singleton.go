package di

import (
	"os"

	"darbelis.eu/taxi/internal/handler"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/internal/web"
	"darbelis.eu/taxi/pkg/nominatim"
	"darbelis.eu/taxi/pkg/util"
	"darbelis.eu/taxi/pkg/ws"
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

var driverStatusChangeForDriverListener *handler.DriverStatusChangeForDriverListener = nil
var driverStatusChangeForInvitationListener *handler.DriverStatusChangeForInvitationListener = nil

var driverStatusChangeForClientListener *handler.DriverStatusChangeForClientListener = nil

var nominatimClient *nominatim.Client = nil

var nominatimAddressWorkerPool *nominatim.AddressWorkerPool = nil

var nominatimCoordinatesWorkerPool *nominatim.CoordinatesWorkerPool = nil

var delegatingDriversMessageHandler *handler.DelegatingMessageHandler = nil

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
		driversHubInstance = ws.NewHub(GetDispatcher(), "driver", messages.Encode, messages.Decode)
	}

	return driversHubInstance
}

func GetPassengersHub() *ws.Hub {
	if passengersHubInstance == nil {
		passengersHubInstance = ws.NewHub(GetDispatcher(), "passenger", messages.Encode, messages.Decode)
	}

	return passengersHubInstance
}

func GetDriversWsHandler() handler.MessageHandler {
	if driversMessagesWsHandler == nil {
		driversMessagesWsHandler = handler.NewDriversMessageHandler(GetMainState(), GetDriversHub(), GetPassengersHub(), GetDispatcher(), GetAddressWorkerPool())
	}

	return driversMessagesWsHandler
}

func GetPassengersWsHandler() handler.MessageHandler {
	if passengersMessagesWsHandler == nil {
		passengersMessagesWsHandler = handler.NewPassengersMessageHandler(GetMainState(), GetDriversHub(), GetPassengersHub(), GetAddressWorkerPool())
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

func GetDriverStatusChangeForDriverListener() *handler.DriverStatusChangeForDriverListener {
	if driverStatusChangeForDriverListener == nil {
		driverStatusChangeForDriverListener = handler.NewDriverStatusChangeForDriverListener(GetMainState())
	}

	return driverStatusChangeForDriverListener
}

func GetDriverStatusChangeForInvitationListener() *handler.DriverStatusChangeForInvitationListener {
	if driverStatusChangeForInvitationListener == nil {
		driverStatusChangeForInvitationListener = handler.NewDriverStatusChangeForInvitationListener(GetMainState(), GetDriversHub(), GetPassengersHub())
	}

	return driverStatusChangeForInvitationListener
}

func GetDriverStatusChangeForClientListener() *handler.DriverStatusChangeForClientListener {

	if driverStatusChangeForClientListener == nil {
		driverStatusChangeForClientListener = handler.NewDriverStatusChangeForClientListener(GetDriversHub())
	}

	return driverStatusChangeForClientListener
}

func GetNominatimClient() *nominatim.Client {
	if nominatimClient == nil {
		server := os.Getenv("NOMINATIM_SERVER")
		user := os.Getenv("NOMINATIM_USER")
		password := os.Getenv("NOMINATIM_PASSWORD")

		nominatimClient = nominatim.NewClient(server, user, password)
	}

	return nominatimClient
}

func GetCoordinatesWorkerPool() *nominatim.CoordinatesWorkerPool {
	if nominatimCoordinatesWorkerPool == nil {
		nominatimCoordinatesWorkerPool = nominatim.NewCoordinatesWorkerPool(GetNominatimClient())
	}

	return nominatimCoordinatesWorkerPool
}

func GetAddressWorkerPool() *nominatim.AddressWorkerPool {
	if nominatimAddressWorkerPool == nil {
		nominatimAddressWorkerPool = nominatim.NewAddressWorkerPool(GetNominatimClient())
	}

	return nominatimAddressWorkerPool
}

func GetDelegatingDriversMessageHandler() *handler.DelegatingMessageHandler {
	if delegatingDriversMessageHandler == nil {
		delegatingDriversMessageHandler = handler.NewDelegatingMessageHandler()
	}

	return delegatingDriversMessageHandler
}
