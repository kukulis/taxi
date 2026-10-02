package web

import (
	"net/http"
	"sort"

	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	"github.com/gin-gonic/gin"
)

// observerTableLimit caps how many drivers/passengers/invitations the observer page shows.
const observerTableLimit = 5

type WebController struct {
	mainState *state.MainState
}

func NewWebController(mainState *state.MainState) *WebController {
	return &WebController{mainState: mainState}
}

func (controller *WebController) Index(c *gin.Context) {
	RollCookie(c)
	c.HTML(http.StatusOK, "main.gohtml", gin.H{})
}
func (controller *WebController) Driver(c *gin.Context) {
	RollCookie(c)
	c.HTML(http.StatusOK, "driver.gohtml", gin.H{})
}
func (controller *WebController) Passenger(c *gin.Context) {
	RollCookie(c)
	c.HTML(http.StatusOK, "passenger.gohtml", gin.H{})
}
func (controller *WebController) Observer(c *gin.Context) {
	RollCookie(c)
	driverFilter := state.NewEmptyFilter()
	driverFilter.Type = string(state.ClientTypeDriver)
	driverClients := controller.mainState.SearchClients(driverFilter)
	drivers := util.ArrayMap(driverClients, func(c state.Client) *state.Driver { return c.(*state.Driver) })
	if len(drivers) > observerTableLimit {
		drivers = drivers[:observerTableLimit]
	}

	passengerFilter := state.NewEmptyFilter()
	passengerFilter.Type = string(state.ClientTypePassenger)
	passengerClients := controller.mainState.SearchClients(passengerFilter)
	passengers := util.ArrayMap(passengerClients, func(c state.Client) *state.Passenger { return c.(*state.Passenger) })
	if len(passengers) > observerTableLimit {
		passengers = passengers[:observerTableLimit]
	}

	// empty filter matches all; the container is a map, so sort for a stable newest-first view
	invitations := controller.mainState.GetInvitationsContainer().GetInvitationsSnapshot(&state.InvitationsFilter{})
	sort.Slice(invitations, func(i, j int) bool {
		return invitations[i].CreatedAt.After(invitations[j].CreatedAt)
	})
	if len(invitations) > observerTableLimit {
		invitations = invitations[:observerTableLimit]
	}

	c.HTML(http.StatusOK, "observer.gohtml", gin.H{
		"Drivers":     drivers,
		"Passengers":  passengers,
		"Invitations": invitations,
	})
}
