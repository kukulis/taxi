// Mirrors internal/dao/driver_search_result.go — keep the two in sync.
export class DriverSearchResult {
    driverId = null;
    driverInfo = null;
    lat = null;
    lon = null;
    distanceKm = null;

    /**
     * @param {Object} obj one driver from the /api/drivers JSON response
     * @returns {DriverSearchResult}
     */
    static fromObject(obj) {
        const result = new DriverSearchResult();
        result.driverId = obj.driver_id;
        result.driverInfo = obj.driver_info;
        result.lat = obj.lat;
        result.lon = obj.lon;
        result.distanceKm = obj.distance_km;

        return result;
    }
}
