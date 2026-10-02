package web

import (
	"net/http"
	"strconv"

	"darbelis.eu/taxi/internal/dto"
	"darbelis.eu/taxi/internal/state"
	"github.com/gin-gonic/gin"
)

type DriverApiController struct {
	mainState *state.MainState
}

func NewDriverApiController(mainState *state.MainState) *DriverApiController {
	return &DriverApiController{mainState: mainState}
}

// Search handles GET /api/driver?searchLat=..&searchLon=..&distance=..
// distance is in kilometers, matching MainState.GetNearestClients' unit.
// Only idle drivers are returned.
func (c *DriverApiController) Search(ctx *gin.Context) {
	searchLat, err := strconv.ParseFloat(ctx.Query("searchLat"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "searchLat is required and must be a number"})
		return
	}

	searchLon, err := strconv.ParseFloat(ctx.Query("searchLon"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "searchLon is required and must be a number"})
		return
	}

	maxDistanceKm, err := strconv.ParseFloat(ctx.Query("distance"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "distance is required and must be a number"})
		return
	}

	// only idle drivers can take a passenger, so only they are offered to be invited
	filter := state.NewEmptyFilter()
	filter.Type = string(state.ClientTypeDriver)
	filter.Active = 1
	filter.DriverStatuses = []state.DriverStatus{state.DriverStatusIdle}

	nearest := c.mainState.GetNearestClients(filter, searchLat, searchLon, maxDistanceKm)

	results := make([]dto.DriverSearchResult, 0, len(nearest))
	for _, clientDistance := range nearest {
		if clientDistance.Distance > maxDistanceKm {
			break // GetNearestClients' box prefilter is square, not circular, and is sorted ascending - so once a result exceeds the max, nothing further will match either
		}
		driver := clientDistance.Client.(*state.Driver)
		results = append(results, dto.DriverSearchResult{
			DriverId:   driver.Id,
			Lat:        driver.Lat,
			Lon:        driver.Lon,
			DistanceKm: clientDistance.Distance,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"drivers": results})
}
