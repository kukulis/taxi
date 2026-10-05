import {Dispatcher} from "./dispatcher.js";
import {ApiClient} from "./api.js";
import {WsService} from "./ws_service.js";
import {DriverComponent} from "./driver_component.js";
import {PassengerComponent} from "./passenger_component.js";
import {CommonComponent} from "./common_component.js";
import {State} from "./state.js";
import {GetE, NewEC, NewECT} from "./util.js";
import {ApplicationEvent, ClientType, InvitationStatus, MessageType, WsEvent} from "./constants.js";

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
     * @param {string} clientType
     * @returns {CommonComponent}
     */
    wireCommonComponent(clientType) {
        const commonComponent = this.getCommonComponent(clientType);
        const dispatcher = this.getDispatcher();
        const wsService = clientType === ClientType.DRIVER ? this.getDriverWs() : this.getPassengerWs();

        dispatcher.addListener(WsEvent.ON_OPEN_WS, () => commonComponent.handleServerRequestsCoordinates({}));
        dispatcher.addListener(WsEvent.ON_CLOSE_WS, () => dispatcher.dispatch(ApplicationEvent.ERROR, {message: 'The WS connection closed'}));
        dispatcher.addListener(ApplicationEvent.SEND_WS_MESSAGE, ([ct, type, data]) => this.sendWsMessage(ct, type, data));

        wsService.addHandle(MessageType.SERVER_SENDS_CLIENT_ADDRESS, (data) => commonComponent.handleServerSendsClientAddress(data));
        wsService.addHandle(MessageType.SERVER_REQUESTS_COORDINATES, (data) => commonComponent.handleServerRequestsCoordinates(data));
        wsService.addHandle(MessageType.SERVER_REQUESTS_CLIENT_INFO, (data) => commonComponent.handleServerRequestsClientInfo(data));

        return commonComponent;
    }

    /**
     * @param {string} driverId
     * @returns {DriverComponent}
     */
    wireDriverComponent(driverId) {
        const driverComponent = this.getDriverComponent(driverId);
        const wsService = this.getDriverWs();

        wsService.addHandle(MessageType.SERVER_OFFERS_PASSENGER, (data) => driverComponent.handleServerOffersPassenger(data));
        wsService.addHandle(MessageType.SERVER_CANCELS_OFFER, (data) => driverComponent.handleServerCancelsOffer(data));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_VOYAGE_STARTED, (data) => driverComponent.handleServerNotifiesVoyageStarted(data));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_VOYAGE_FINISHED, (data) => driverComponent.handleServerNotifiesVoyageFinished(data));
        wsService.addHandle(MessageType.SERVER_REFRESH_DRIVER_STATUS, (data) => driverComponent.handleServerRefreshStatus(data));
        wsService.addHandle(MessageType.SERVER_REFRESH_DRIVER_OFFERS, (data) => driverComponent.handleServerRefreshOffers(data));

        return driverComponent;
    }

    /**
     * @param {string} passengerId
     * @returns {PassengerComponent}
     */
    wirePassengerComponent(passengerId) {
        const passengerComponent = this.getPassengerComponent(passengerId);
        const wsService = this.getPassengerWs();

        wsService.addHandle(MessageType.SERVER_NOTIFIES_INVITE_ACCEPTED, (data) => passengerComponent.updateInvitationStatus(data.driver_id, InvitationStatus.ACCEPTED));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_INVITE_REJECTED, (data) => passengerComponent.updateInvitationStatus(data.driver_id, InvitationStatus.REJECTED));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_INVITE_CANCELED, (data) => passengerComponent.updateInvitationStatus(data.driver_id, InvitationStatus.CANCELED));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_VOYAGE_STARTED, (data) => passengerComponent.updateInvitationStatus(data.driver_id, InvitationStatus.DRIVING));
        wsService.addHandle(MessageType.SERVER_NOTIFIES_VOYAGE_FINISHED, (data) => passengerComponent.updateInvitationStatus(data.driver_id, InvitationStatus.COMPLETED));
        wsService.addHandle(MessageType.SERVER_REFRESH_PASSENGER_INVITATIONS, (data) => passengerComponent.handleServerRefreshInvitations(data));

        return passengerComponent;
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
