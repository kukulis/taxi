export class Offer {
    lat = null;
    lng = null;
    passengerId = null;
    status = null;
    time = null;

    setLat(lat) {
        this.lat = lat;
        return this;
    }

    setLng(lng) {
        this.lng = lng;
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
        return this.lng;
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
}
