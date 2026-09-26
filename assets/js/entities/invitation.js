export class Invitation {
    lat = null;
    lon = null;
    driverId = null;
    status = null;
    time = null;

    setLat(lat) {
        this.lat = lat;
        return this;
    }

    setLon(lng) {
        this.lon = lng;
        return this;
    }

    setDriverId(driverId) {
        this.driverId = driverId;
        return this;
    }

    setStatus(status) {
        this.status = status;
        return this;
    }

    setTime(time) {
        this.time = time;
        return this;
    }

    getLat() {
        return this.lat;
    }

    getLon() {
        return this.lon;
    }

    getDriverId() {
        return this.driverId;
    }

    getStatus() {
        return this.status;
    }

    getTime() {
        return this.time;
    }

    fromObject(obj) {
        this.lat = obj.lat;
        this.lon = obj.lon;
        this.driverId = obj.driver_id;
        this.status = obj.status;
        this.time = obj.time;
        return this;
    }
}
