// Shared application state, alive for the page's whole lifetime -
// the frontend counterpart to the backend's MainState.
export class State {

    lat = null;
    lon = null;
    address = "";
    phone = "";
    name = "";

    setCoordinates(lat, lon) {
        this.lat = lat;
        this.lon = lon;
    }

    getCoordinates() {
        return {lat: this.lat, lon: this.lon};
    }

    setAddress(address) {
        this.address = address;
    }

    getAddress() {
        return this.address;
    }
}
