import {DriverSearchResult} from "./entities/driver_search_result.js";

export class ApiClient {
    dispatcher = null;

    constructor(dispatcher) {
        this.dispatcher = dispatcher;
    }

    /**
     * @param {number} searchLat
     * @param {number} searchLon
     * @param {number} distanceKm max distance from the search point, in kilometers
     * @returns {Promise<DriverSearchResult[]>} sorted by distance, nearest first
     */
    async getDrivers(searchLat, searchLon, distanceKm) {
        try {
            const params = new URLSearchParams({searchLat, searchLon, distance: distanceKm});
            const response = await fetch('/api/driver?' + params);
            // throws on a non-JSON body (e.g. a proxy's HTML error page), handled below
            const data = await response.json();

            if (!response.ok) {
                this.dispatcher.dispatch('error', {message: data.error});
                return [];
            }

            return data.drivers.map( (element) => new DriverSearchResult().fromObject(element));
        } catch (e) {
            // network failure or unparseable response
            this.dispatcher.dispatch('error', {message: e.message});
            return [];
        }
    }
}
