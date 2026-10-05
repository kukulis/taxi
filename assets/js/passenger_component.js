import {AppChildren, ClearE, NewEC, NewECT, NewT} from "./util.js";
import {ApplicationEvent, ClientType, InvitationStatus, MessageType} from "./constants.js";
import {DriverSearchResult} from "./entities/driver_search_result.js";
import {Invitation} from "./entities/invitation.js";
import {ServerRefreshPassengerInvitations} from "./entities/messages.js";

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

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {State}
     */
    state = null;

    /**
     * @type {DriverSearchResult[]}
     */
    driversSearchResults = [];

    /**
     *
     * @type {[Invitation]}
     */
    invitations = [];

    constructor(passengerId, apiClient, dispatcher, state) {
        this.passengerId = passengerId;
        this.apiClient = apiClient;
        this.dispatcher = dispatcher;
        this.state = state;
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
        const {lat, lon} = this.state.getCoordinates();
        if (lat === null || lon === null) {
            console.warn('no passenger coordinates yet, cannot search drivers');
            return;
        }

        this.driversSearchResults = await this.apiClient.getDrivers(lat, lon, distanceKm);
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
        const {lat, lon} = this.state.getCoordinates();
        if (lat === null || lon === null) {
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
            lat: lat,
            lon: lon,
        });

        // the list goes stale quickly (drivers change status), so drop it; a distance button reloads it
        this.driversSearchResults = [];
        this.renderDrivers();
    }

    renderInfo() {
        ClearE(this.infoView);
        this.infoView.appendChild(NewECT('h3', 'passenger-title', 'Passenger'));
        this.infoView.appendChild(NewECT('div', 'passenger-id-info', this.passengerId));
    }

    handleServerRefreshInvitations(data) {
        const message = new ServerRefreshPassengerInvitations().fromObject(data);

        this.invitations = message.getInvitations();
        this.renderInvitations();
    }

    sendMessage(type, data) {
        console.log('passenger component dispatches send ws message', type, data);
        this.dispatcher.dispatch(ApplicationEvent.SEND_WS_MESSAGE, [ClientType.PASSENGER, type, data]);
    }
}
