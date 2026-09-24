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
    fromObject(obj) {
        // const result = new DriverSearchResult();
        this.driverId = obj.driver_id;
        this.driverInfo = obj.driver_info;
        this.lat = obj.lat;
        this.lon = obj.lon;
        this.distanceKm = obj.distance_km;

        return this;
    }

    setDriverId(driverId) {
        this.driverId = driverId;
        return this;
    }

    setDriverInfo(driverInfo) {
        this.driverInfo = driverInfo;
        return this;
    }

    setLat(lat) {
        this.lat = lat;
        return this;
    }

    setLon(lon) {
        this.lon = lon;
        return this;
    }

    setDistanceKm(distanceKm) {
        this.distanceKm = distanceKm;
        return this;
    }
}
