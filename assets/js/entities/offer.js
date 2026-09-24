export class Offer {
    lat = null;
    lon = null;
    /**
     * @type {string}
     */
    passengerId = null;
    status = null;
    time = null;

    setLat(lat) {
        this.lat = lat;
        return this;
    }

    setLng(lng) {
        this.lon = lng;
        return this;
    }

    setPassengerId(passengerId) {
        this.passengerId = passengerId;
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

    getLng() {
        return this.lon;
    }

    getPassengerId() {
        return this.passengerId;
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
        this.passengerId = obj.passenger_id;
        this.status = obj.status;
        this.time = obj.time;
        return this;
    }
}
