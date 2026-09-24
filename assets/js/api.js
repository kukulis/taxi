export class ApiClient {
    getDrivers() {
        return fetch('/api/drivers')
            .then(response => response.json())
            .then(data => data.drivers);
    }
}