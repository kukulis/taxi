package web

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"darbelis.eu/taxi/internal/dto"
	"darbelis.eu/taxi/internal/state"
	"darbelis.eu/taxi/pkg/util"
	"github.com/gin-gonic/gin"
)

func TestDriverApiController_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name        string
		drivers     []*state.Driver
		query       string
		wantStatus  int
		wantError   string
		wantDrivers []dto.DriverSearchResult
	}{
		{
			name:       "missing searchLat returns 400",
			query:      "?searchLon=25.2797&distance=5",
			wantStatus: http.StatusBadRequest,
			wantError:  "searchLat is required and must be a number",
		},
		{
			name:       "missing searchLon returns 400",
			query:      "?searchLat=54.6872&distance=5",
			wantStatus: http.StatusBadRequest,
			wantError:  "searchLon is required and must be a number",
		},
		{
			name:       "missing distance returns 400",
			query:      "?searchLat=54.6872&searchLon=25.2797",
			wantStatus: http.StatusBadRequest,
			wantError:  "distance is required and must be a number",
		},
		{
			name: "idle driver within distance is returned",
			drivers: []*state.Driver{
				{Id: "near-driver", Lat: 54.6892, Lon: 25.2799, Status: state.DriverStatusIdle}, // ~0.22km from Vilnius center
			},
			query:      "?searchLat=54.6872&searchLon=25.2797&distance=5",
			wantStatus: http.StatusOK,
			wantDrivers: []dto.DriverSearchResult{
				{DriverId: "near-driver", Lat: 54.6892, Lon: 25.2799, DistanceKm: 0.2227610609715733},
			},
		},
		{
			name: "non-idle driver is excluded even if nearby",
			drivers: []*state.Driver{
				{Id: "working-driver", Lat: 54.6892, Lon: 25.2799, Status: state.DriverStatusWorking}, // ~0.22km from Vilnius center
			},
			query:       "?searchLat=54.6872&searchLon=25.2797&distance=5",
			wantStatus:  http.StatusOK,
			wantDrivers: []dto.DriverSearchResult{},
		},
		{
			name: "driver beyond max distance is excluded",
			drivers: []*state.Driver{
				{Id: "far-driver", Lat: 54.8985, Lon: 23.9036, Status: state.DriverStatusIdle}, // ~91.29km from Vilnius center (Kaunas)
			},
			query:       "?searchLat=54.6872&searchLon=25.2797&distance=5",
			wantStatus:  http.StatusOK,
			wantDrivers: []dto.DriverSearchResult{},
		},
		{
			name: "two idle drivers sorted ascending by distance",
			drivers: []*state.Driver{
				{Id: "middle-driver", Lat: 54.7500, Lon: 25.0000, Status: state.DriverStatusIdle}, // ~19.27km from Vilnius center
				{Id: "near-driver", Lat: 54.6892, Lon: 25.2799, Status: state.DriverStatusIdle},   // ~0.22km from Vilnius center
			},
			query:      "?searchLat=54.6872&searchLon=25.2797&distance=50",
			wantStatus: http.StatusOK,
			wantDrivers: []dto.DriverSearchResult{
				{DriverId: "near-driver", Lat: 54.6892, Lon: 25.2799, DistanceKm: 0.2227610609715733},
				{DriverId: "middle-driver", Lat: 54.7500, Lon: 25.0000, DistanceKm: 19.273340171657384},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mainState := state.NewMainState(util.NewFixedClock(time.Now()))
			for _, driver := range c.drivers {
				mainState.AddClient(driver)
			}

			controller := NewDriverApiController(mainState)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/api/driver"+c.query, nil)

			controller.Search(ctx)

			if recorder.Code != c.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", recorder.Code, c.wantStatus, recorder.Body.String())
			}

			if c.wantError != "" {
				var errBody struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(recorder.Body.Bytes(), &errBody); err != nil {
					t.Fatalf("failed to decode error body: %v, body = %s", err, recorder.Body.String())
				}
				if errBody.Error != c.wantError {
					t.Errorf("error = %q, want %q", errBody.Error, c.wantError)
				}
				return
			}

			var okBody struct {
				Drivers []dto.DriverSearchResult `json:"drivers"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &okBody); err != nil {
				t.Fatalf("failed to decode success body: %v, body = %s", err, recorder.Body.String())
			}

			if len(okBody.Drivers) != len(c.wantDrivers) {
				t.Fatalf("drivers = %+v, want %+v", okBody.Drivers, c.wantDrivers)
			}
			for i, want := range c.wantDrivers {
				got := okBody.Drivers[i]
				if got.DriverId != want.DriverId || got.Lat != want.Lat || got.Lon != want.Lon {
					t.Errorf("drivers[%d] = %+v, want %+v", i, got, want)
				}
				if math.Abs(got.DistanceKm-want.DistanceKm) > 1e-9 {
					t.Errorf("drivers[%d].DistanceKm = %v, want %v", i, got.DistanceKm, want.DistanceKm)
				}
			}
		})
	}
}

// A passenger sitting at the exact same coordinates as a driver must never
// show up in a driver search, even though both are in the same MainState.clients
// map and SearchClients' box prefilter alone can't tell them apart by position.
func TestDriverApiController_Search_ExcludesPassengerAtSameCoordinates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const lat, lon = 54.6892, 25.2799

	mainState := state.NewMainState(util.NewFixedClock(time.Now()))
	mainState.AddClient(&state.Driver{Id: "driver-1", Lat: lat, Lon: lon, Status: state.DriverStatusIdle})
	mainState.AddClient(&state.Passenger{Id: "passenger-1", Lat: lat, Lon: lon})

	controller := NewDriverApiController(mainState)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/driver?searchLat=54.6892&searchLon=25.2799&distance=5", nil)

	controller.Search(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var body struct {
		Drivers []dto.DriverSearchResult `json:"drivers"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v, body = %s", err, recorder.Body.String())
	}

	if len(body.Drivers) != 1 {
		t.Fatalf("drivers = %+v, want exactly 1 (the driver, not the passenger)", body.Drivers)
	}
	if body.Drivers[0].DriverId != "driver-1" {
		t.Errorf("drivers[0].DriverId = %q, want %q", body.Drivers[0].DriverId, "driver-1")
	}
}
