package dao

// DriverSearchResult is the JSON shape for one driver in a nearest-drivers search result.
type DriverSearchResult struct {
	DriverId   string  `json:"driver_id"`
	DriverInfo string  `json:"driver_info"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	DistanceKm float64 `json:"distance_km"`
}
