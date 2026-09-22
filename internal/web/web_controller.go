package web

import (
	"net/http"

	"darbelis.eu/taxi/internal/state"
	"github.com/gin-gonic/gin"
)

// observerTableLimit caps how many drivers/passengers the observer page shows.
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
	drivers := controller.mainState.GetDriversSnapshot()
	if len(drivers) > observerTableLimit {
		drivers = drivers[:observerTableLimit]
	}

	passengers := controller.mainState.GetPassengersSnapshot()
	if len(passengers) > observerTableLimit {
		passengers = passengers[:observerTableLimit]
	}

	c.HTML(http.StatusOK, "observer.gohtml", gin.H{
		"Drivers":    drivers,
		"Passengers": passengers,
	})
}
