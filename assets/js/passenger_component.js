import {AppChildren, ClearE, NewEC, NewECT, NewT} from "./util.js";
import {InvitationStatus, MessageType} from "./constants.js";
import {DriverSearchResult} from "./entities/driver_search_result.js";
import {Invitation} from "./entities/invitation.js";

export class PassengerComponent {

    passengerId = null;
    // main view
    passengerView = null;
    infoView = null;
    // invitations sent to drivers
    invitationsView = null;
    // distance buttons + drivers list
    driversView = null;
    driversListView = null;
    // websocket connection
    conn = null;

    lat = null;
    lon = null;

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {DriverSearchResult[]}
     */
    driversSearchResults = [];

    /**
     *
     * @type {[Invitation]}
     */
    invitations = [
        // (new Invitation()).setDriverId('1234567890')
        //     .setStatus(InvitationStatus.PENDING)
        //     .setLat(12)
        //     .setLon(50)
    ];

    constructor(passengerId, apiClient) {
        this.passengerId = passengerId;
        this.apiClient = apiClient;
    }

    async render() {
        this.passengerView = NewEC('div', 'passenger-component');
        this.infoView = NewEC('div', 'passenger-info');
        this.invitationsView = NewEC('div', 'passenger-invitations');
        this.driversView = NewEC('div', 'passenger-drivers');

        this.renderInfo();
        this.renderInvitations();
        this.renderDriversView();

        return AppChildren(this.passengerView, [this.infoView, this.invitationsView, this.driversView]);
    }

    /**
     * Status text with a per-status class (invitation-status-pending, ...) for its color.
     * @param {string} className
     * @param {string} status
     * @returns {HTMLElement}
     */
    statusBadge(className, status) {
        const badge = NewECT('span', className, status);
        badge.classList.add('invitation-status-' + status);
        return badge;
    }

    renderInvitations() {
        ClearE(this.invitationsView);

        this.invitationsView.appendChild(NewECT('h3', 'invitations-title', 'Invitations'));

        for (const invitation of this.invitations) {
            const parts = [
                NewECT('span', 'invitation-driver', invitation.getDriverId()),
                NewECT('span', 'invitation-coords', `(${invitation.getLat()}, ${invitation.getLon()})`),
                this.statusBadge('invitation-status', invitation.getStatus()),
            ];

            const actions = NewEC('span', 'invitation-actions');
            if (this.isActiveInvitation(invitation)) {
                const cancelButton = NewECT('button', 'invitation-cancel', 'cancel');
                cancelButton.addEventListener('click', () => this.onInvitationCancelClick(invitation));
                actions.appendChild(cancelButton);
            } else if (invitation.getStatus() === InvitationStatus.REJECTED
                || invitation.getStatus() === InvitationStatus.COMPLETED) {
                const hideButton = NewECT('button', 'invitation-hide', 'hide');
                hideButton.addEventListener('click', () => this.onInvitationHideClick(invitation));
                actions.appendChild(hideButton);
            }
            if (actions.hasChildNodes()) {
                parts.push(actions);
            }

            const row = NewEC('div', 'invitation-row');
            parts.forEach((part, i) => {
                if (i > 0) {
                    row.appendChild(NewECT('span', 'invitation-divider', '|'));
                }
                row.appendChild(part);
            });

            this.invitationsView.appendChild(row);
        }
    }

    /**
     * Pending or accepted: the passenger can still cancel it and the server may still update it.
     * @param {Invitation} invitation
     * @returns {boolean}
     */
    isActiveInvitation(invitation) {
        return invitation.getStatus() === InvitationStatus.PENDING
            || invitation.getStatus() === InvitationStatus.ACCEPTED;
    }

    /**
     * @param {Invitation} invitation
     */
    onInvitationCancelClick(invitation) {
        this.sendMessage(MessageType.PASSENGER_CANCELS_INVITE, {driver_id: invitation.getDriverId()});

        this.invitations = this.invitations.filter((i) => i !== invitation);
        this.renderInvitations();
    }

    /**
     * Removes a finished invitation from the list; the server keeps it.
     * @param {Invitation} invitation
     */
    onInvitationHideClick(invitation) {
        this.invitations = this.invitations.filter((i) => i !== invitation);
        this.renderInvitations();
    }

    /**
     * Server notifications only carry driver_id, so they apply to the latest active invitation to that driver.
     * @param {string} driverId
     * @param {string} status
     */
    updateInvitationStatus(driverId, status) {
        // driving is not cancelable (so not "active"), but the server still finishes it
        const invitation = this.invitations.findLast(
            (i) => i.getDriverId() === driverId
                && (this.isActiveInvitation(i) || i.getStatus() === InvitationStatus.DRIVING),
        );

        if (!invitation) {
            console.warn('no open invitation for driver', driverId, 'to set status', status);
            return;
        }

        invitation.setStatus(status);
        this.renderInvitations();
    }

    renderDriversView() {
        ClearE(this.driversView);

        this.driversView.appendChild(NewECT('h3', 'drivers-title', 'Drivers'));

        const buttons = NewEC('div', 'distance-buttons');
        for (const [label, distanceKm] of [['250m', 0.25], ['500m', 0.5], ['1km', 1], ['5km', 5]]) {
            const button = NewECT('button', 'distance-button', label);
            button.addEventListener('click', () => this.onDistanceClick(distanceKm));
            buttons.appendChild(button);
        }

        this.driversListView = NewEC('div', 'drivers-list');

        AppChildren(this.driversView, [buttons, this.driversListView]);

        this.renderDrivers();
    }

    /**
     * @param {number} distanceKm
     */
    async onDistanceClick(distanceKm) {
        if (this.lat === null || this.lon === null) {
            console.warn('no passenger coordinates yet, cannot search drivers');
            return;
        }

        this.driversSearchResults = await this.apiClient.getDrivers(this.lat, this.lon, distanceKm);
        this.renderDrivers();
    }

    renderDrivers() {
        ClearE(this.driversListView);

        for (const driver of this.driversSearchResults) {
            const parts = [
                NewECT('span', 'driver-row-id', driver.driverId),
                NewECT('span', 'driver-row-info', driver.driverInfo),
                NewECT('span', 'driver-row-coords', `(${driver.lat}, ${driver.lon})`),
                NewECT('span', 'driver-row-distance', `${driver.distanceKm.toFixed(2)} km`),
            ];

            const inviteButton = NewECT('button', 'driver-row-invite', 'invite');
            inviteButton.addEventListener('click', () => this.onInviteClick(driver));
            parts.push(inviteButton);

            const row = NewEC('div', 'driver-row');
            parts.forEach((part, i) => {
                if (i > 0) {
                    row.appendChild(NewECT('span', 'driver-row-divider', '|'));
                }
                row.appendChild(part);
            });

            this.driversListView.appendChild(row);
        }
    }

    /**
     * @param {DriverSearchResult} driver
     */
    onInviteClick(driver) {
        if (this.lat === null || this.lon === null) {
            console.warn('no passenger coordinates yet, cannot invite a driver');
            return;
        }

        // the invitation row shows where the driver is; the message carries where the passenger is (pickup point)
        this.invitations.push(
            new Invitation()
                .setDriverId(driver.driverId)
                .setLat(driver.lat)
                .setLon(driver.lon)
                .setStatus(InvitationStatus.PENDING)
                .setTime(new Date()),
        );
        this.renderInvitations();

        this.sendMessage(MessageType.PASSENGER_INVITES_DRIVER, {
            driver_id: driver.driverId,
            lat: this.lat,
            lon: this.lon,
        });
    }

    renderInfo() {
        ClearE(this.infoView);
        this.infoView.appendChild(NewECT('h3', 'passenger-title', 'Passenger'));
        this.infoView.appendChild(NewECT('div', 'passenger-id-info', this.passengerId));

        if (this.lat !== null && this.lon !== null) {
            this.infoView.appendChild(NewECT('div', 'passenger-location', `(${this.lat}, ${this.lon})`));
        } else {
            this.infoView.appendChild(NewECT('div', 'passenger-location', '...'));
        }
    }

    initWs(isTls) {
        let protocol = 'ws://';

        if (isTls) {
            protocol = 'wss://';
        }

        if (!window["WebSocket"]) {
            this.passengerView.appendChild(NewT('Your browser does not support WebSockets.'));
            return;
        }
        this.conn = new WebSocket(protocol + document.location.host + "/ws-passenger");
        // send our coordinates right away instead of waiting for the server to ask;
        // only once open, since send() throws while the socket is still connecting
        this.conn.onopen = () => {
            this.handleServerRequestsCoordinates({});
        };
        this.conn.onclose = (evt) => {
            this.passengerView.appendChild(NewT('Websocket connection closed.'));
        };
        this.conn.onmessage = (evt) => {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                this.handleWsMessage(messages[i]);
            }
        };
    }

    handleWsMessage(message) {
        const envelope = JSON.parse(message);

        switch (envelope.type) {
            case MessageType.SERVER_REQUESTS_COORDINATES:
                this.handleServerRequestsCoordinates(envelope.data);
                break;
            case MessageType.SERVER_REQUESTS_CLIENT_INFO:
                this.handleServerRequestsClientInfo(envelope.data);
                break;
            case MessageType.SERVER_NOTIFIES_INVITE_ACCEPTED:
                this.updateInvitationStatus(envelope.data.driver_id, InvitationStatus.ACCEPTED);
                break;
            case MessageType.SERVER_NOTIFIES_INVITE_REJECTED:
                this.updateInvitationStatus(envelope.data.driver_id, InvitationStatus.REJECTED);
                break;
            case MessageType.SERVER_NOTIFIES_INVITE_CANCELED:
                this.updateInvitationStatus(envelope.data.driver_id, InvitationStatus.CANCELED);
                break;
            case MessageType.SERVER_NOTIFIES_VOYAGE_STARTED:
                this.updateInvitationStatus(envelope.data.driver_id, InvitationStatus.DRIVING);
                break;
            case MessageType.SERVER_NOTIFIES_VOYAGE_FINISHED:
                this.updateInvitationStatus(envelope.data.driver_id, InvitationStatus.COMPLETED);
                break;
            default:
                console.log('unhandled message type', envelope.type, envelope);
        }
    }

    handleServerRequestsCoordinates(data) {
        console.log('server requested coordinates');

        if (!navigator.geolocation) {
            console.error('geolocation is not supported by this browser');
            this.sendMessage(MessageType.CLIENT_ERROR, {error: 'Geolocation is not supported by this browser.'});
            return;
        }

        navigator.geolocation.getCurrentPosition(
            (position) => {
                this.lat = position.coords.latitude;
                this.lon = position.coords.longitude;
                console.log('got coordinates', this.lat, this.lon);
                this.sendMessage(MessageType.CLIENT_RESPONDS_COORDINATES, {lat: this.lat, lon: this.lon});
                this.renderInfo();
            },
            (error) => {
                console.error('failed to get coordinates:', error.message);
                this.sendMessage(MessageType.CLIENT_ERROR, {error: error.message});
            },
            {enableHighAccuracy: true, maximumAge: 5000, timeout: 10000}
        );
    }

    handleServerRequestsClientInfo(data) {
        // TODO: real phone once there's an input for it; empty for now.
        console.log('server requested client info, replying with empty values for now');
        this.sendMessage(MessageType.CLIENT_RESPONDS_INFO, {});
    }

    sendMessage(type, data) {
        console.log('sending message', type, data);
        this.conn.send(JSON.stringify({
            id: crypto.randomUUID(),
            type: type,
            data: data,
        }));
    }
}
