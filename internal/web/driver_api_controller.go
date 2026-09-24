package web

import (
	"net/http"
	"strconv"

	"darbelis.eu/taxi/internal/dao"
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
// distance is in kilometers, matching MainState.GetNearestDrivers' unit.
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

	nearest := c.mainState.GetNearestDrivers(searchLat, searchLon)

	results := make([]dao.DriverSearchResult, 0, len(nearest))
	for _, driverDistance := range nearest {
		if driverDistance.Distance > maxDistanceKm {
			break // GetNearestDrivers is sorted ascending, so nothing further will match either
		}
		results = append(results, dao.DriverSearchResult{
			DriverId:   driverDistance.Driver.Id,
			Lat:        driverDistance.Driver.Lat,
			Lon:        driverDistance.Driver.Lon,
			DistanceKm: driverDistance.Distance,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"drivers": results})
}
