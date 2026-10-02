package state

import (
	"testing"
	"time"

	"darbelis.eu/taxi/pkg/test_utils"
	"darbelis.eu/taxi/pkg/util"
)

func TestMainState_GetNearestClients(t *testing.T) {
	cases := []nearestClientTestCase{
		provideMixedDriverAndPassenger(),
		provideInactiveClientExcluded(),
		provideNoCoordinatesClient(),
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.mainState.GetNearestClients(c.clientFilter, c.centerLat, c.centerLon, c.distanceKm)

			ids := util.ArrayMap(got, func(cd ClientDistance) string { return cd.Client.getId() })
			distances := util.ArrayMap(got, func(cd ClientDistance) float64 { return cd.RoundDistance().Distance })

			test_utils.AssertEquals(c.expectedIds, ids, t, false, c.name+" ids")
			test_utils.AssertEquals(c.expectedDistances, distances, t, false, c.name+" distances")
		})
	}
}

// ========= test case structure ===========================================

type nearestClientTestCase struct {
	name              string
	mainState         *MainState
	clientFilter      *ClientFilter
	centerLat         float64
	centerLon         float64
	distanceKm        float64
	expectedIds       []string
	expectedDistances []float64
}

// ========== test cases ========================================================

func provideMixedDriverAndPassenger() nearestClientTestCase {
	const vilniusLat, vilniusLon = 54.6872, 25.2797
	s := NewMainState(util.NewFixedClock(time.Now()))
	s.AddClient(&Driver{Id: "near-driver", Lat: 54.6892, Lon: 25.2799, Active: true})      // ~0.22km from Vilnius center
	s.AddClient(&Driver{Id: "middle-passenger", Lat: 54.7500, Lon: 25.0000, Active: true}) // ~19.27km from Vilnius center
	s.AddClient(&Driver{Id: "far-driver", Lat: 54.8985, Lon: 23.9036, Active: true})       // ~91.29km from Vilnius center (Kaunas) — outside the 50km box below
	f := NewEmptyFilter()
	f.Active = 1

	return nearestClientTestCase{
		name:              "mixed driver and passenger, sorted ascending, far one excluded",
		mainState:         s,
		clientFilter:      f,
		centerLat:         vilniusLat,
		centerLon:         vilniusLon,
		distanceKm:        50,
		expectedIds:       []string{"near-driver", "middle-passenger"},
		expectedDistances: []float64{0.22, 19.27},
	}
}

func provideInactiveClientExcluded() nearestClientTestCase {
	const vilniusLat, vilniusLon = 54.6872, 25.2797
	s := NewMainState(util.NewFixedClock(time.Now()))

	s.AddClient(&Driver{Id: "active-driver", Lat: 54.6892, Lon: 25.2799}) // ~0.22km from Vilnius center

	inactiveDriver := &Driver{Id: "inactive-driver", Lat: 54.6895, Lon: 25.2800} // ~0.26km from Vilnius center
	s.AddClient(inactiveDriver)
	// TODO Letting modify the driver object inside the MainState may be the source of future bugs
	inactiveDriver.Active = false // AddClient forces Active=true; set back to false after

	return nearestClientTestCase{
		name:              "inactive client excluded when filter requires active",
		mainState:         s,
		clientFilter:      &ClientFilter{Active: 1, Online: -1},
		centerLat:         vilniusLat,
		centerLon:         vilniusLon,
		distanceKm:        5,
		expectedIds:       []string{"active-driver"},
		expectedDistances: []float64{0.22},
	}
}

func provideNoCoordinatesClient() nearestClientTestCase {
	const vilniusLat, vilniusLon = 54.6872, 25.2797
	s := NewMainState(util.NewFixedClock(time.Now()))

	s.AddClient(&Driver{Id: "has-coords", Lat: 54.6892, Lon: 25.2799}) // ~0.22km from Vilnius center
	s.AddClient(&Driver{Id: "no-coords"})                              // Lat/Lon left at zero value — never reported coordinates

	return nearestClientTestCase{
		name:              "client with no coordinates stays at (0,0), falls outside the box",
		mainState:         s,
		clientFilter:      NewEmptyFilter(),
		centerLat:         vilniusLat,
		centerLon:         vilniusLon,
		distanceKm:        5,
		expectedIds:       []string{"has-coords"},
		expectedDistances: []float64{0.22},
	}
}
