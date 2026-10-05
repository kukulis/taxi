import {Dispatcher} from "./dispatcher.js";
import {ApiClient} from "./api.js";
import {WsService} from "./ws_service.js";
import {DriverComponent} from "./driver_component.js";
import {PassengerComponent} from "./passenger_component.js";
import {CommonComponent} from "./common_component.js";
import {State} from "./state.js";
import {GetE, NewEC, NewECT} from "./util.js";
import {ApplicationEvent, ClientType} from "./constants.js";

// Holds page-wide services; each one is created on first use.
export class App {
    /** @type {Dispatcher|null} */
    dispatcher = null;
    /** @type {ApiClient|null} */
    apiClient = null;
    /** @type {WsService|null} */
    driverWs = null;
    /** @type {WsService|null} */
    passengerWs = null;
    /** @type {DriverComponent|null} */
    driverComponent = null;
    /** @type {PassengerComponent|null} */
    passengerComponent = null;
    /** @type {CommonComponent|null} */
    commonComponent = null;
    /** @type {State|null} */
    state = null;

    /**
     * @returns {Dispatcher}
     */
    getDispatcher() {
        if (this.dispatcher === null) {
            this.dispatcher = new Dispatcher();
            this.dispatcher.addListener(ApplicationEvent.ERROR, (error) => this.showError(error));
        }

        return this.dispatcher;
    }

    /**
     * @returns {State}
     */
    getState() {
        if (this.state === null) {
            this.state = new State();
        }

        return this.state;
    }

    /**
     * @returns {ApiClient}
     */
    getApiClient() {
        if (this.apiClient === null) {
            this.apiClient = new ApiClient(this.getDispatcher());
        }

        return this.apiClient;
    }

    /**
     * @returns {WsService}
     */
    getDriverWs() {

        if (this.driverWs === null) {
            const isTls = document.location.protocol === 'https:';
            this.driverWs = new WsService('/ws-driver', isTls, this.getDispatcher());
        }

        return this.driverWs;
    }

    /**
     * @returns {WsService}
     */
    getPassengerWs() {
        if (this.passengerWs === null) {
            this.passengerWs = new WsService('/ws-passenger', document.location.protocol === 'https:', this.getDispatcher());
        }

        return this.passengerWs;
    }

    /**
     * @param {string} driverId
     * @returns {DriverComponent}
     */
    getDriverComponent(driverId) {
        if (this.driverComponent === null) {
            this.driverComponent = new DriverComponent(driverId, this.getDispatcher());
        }

        return this.driverComponent;
    }

    /**
     * @param {string} passengerId
     * @returns {PassengerComponent}
     */
    getPassengerComponent(passengerId) {
        if (this.passengerComponent === null) {
            this.passengerComponent = new PassengerComponent(passengerId, this.getApiClient(), this.getDispatcher(), this.getState());
        }

        return this.passengerComponent;
    }

    /**
     * @param {string} clientType
     * @returns {CommonComponent}
     */
    getCommonComponent(clientType) {
        if (this.commonComponent === null) {
            this.commonComponent = new CommonComponent(clientType, this.getDispatcher(), this.getState());
        }

        return this.commonComponent;
    }

    /**
     * Appends a closable error message to the page's #errors div.
     * @param {{message: string}} error
     */
    showError(error) {
        console.error("error", error);

        const errors = GetE("errors");
        if (errors === null) {
            console.error("errors div not found");
            return;
        }

        const errorDiv = NewEC('div', 'error');
        errorDiv.appendChild(NewECT('span', 'error-message', error.message));

        const closeButton = NewECT('button', 'error-close', '✕');
        closeButton.addEventListener('click', () => errorDiv.remove());
        errorDiv.appendChild(closeButton);

        errors.appendChild(errorDiv);
    }

    sendWsMessage (clientType, messageType, data ) {
        if (clientType === ClientType.DRIVER ) {
            this.getDriverWs().sendMessage(messageType, data)
        }
        else if (clientType === ClientType.PASSENGER) {
            this.getPassengerWs().sendMessage(messageType, data)
        }
        else {
            this.getDispatcher().dispatch(ApplicationEvent.ERROR, {message: 'Sending ws message, unknown client type '+clientType})
        }
    }

}
