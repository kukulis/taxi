import {AppChildren, ClearE, NewEC, NewECT, NewT} from "./util.js";
import {DriverStatus, InvitationStatus, MessageType} from "./constants.js";
import {Offer} from "./entities/offer.js";

export class DriverComponent {

    driverId = null;
    // main view
    driverView = null;
    // status squares container
    statusView = null;
    // offers table container
    offersView = null;
    // websocket connection
    conn = null;

    currentStatus =  DriverStatus.RESTING;
    // status square the user clicked on, but has not confirmed with 'select' yet
    selectedStatus = null;

    /**
     *
     * @type {Offer[]}
     */
    offers = [];

    constructor(driverId) {
        this.driverId = driverId;


        // TODO remove after test
        // TODO load offers from REST api on a first load. This will require an API endpoint in the server too.
        this.offers.push(
            new Offer()
                .setLat(12)
                .setLng(50)
                .setStatus(InvitationStatus.PENDING)
                .setPassengerId('123123')
                .setTime(new Date()),
        );
    }

    async render() {
        this.driverView = NewEC('div', 'driver-component');
        this.statusView = NewEC('div', 'driver-statuses');
        this.offersView = NewEC('div', 'driver-offers');

        this.renderStatuses();
        this.renderOffers();

        return AppChildren(this.driverView, [this.statusView, this.offersView]);
    }

    renderStatuses() {
        ClearE(this.statusView);

        this.statusView.appendChild(NewECT('h3', 'driver-title', 'Driver'));
        this.statusView.appendChild(NewECT( 'div', 'driver-id-info',  this.driverId));

        const squaresDiv = NewEC('div', 'status-squares-container');
        for (const status of [DriverStatus.IDLE, DriverStatus.WORKING, DriverStatus.RESTING]) {
            const square = NewEC('div', 'status-square');
            square.classList.add('status-' + status);
            square.appendChild(NewECT('div', 'status-label', status));

            if (status === this.currentStatus) {
                square.appendChild(NewECT('div', 'status-check', '✓'));
            }

            if (status === this.selectedStatus) {
                const button = NewECT('button', 'status-select', 'select');
                button.addEventListener('click', (event) => {
                    // don't let the square's click handler re-render under us
                    event.stopPropagation();
                    this.onStatusSelectClick(status);
                });
                square.appendChild(button);
            }

            square.addEventListener('click', () => {
                this.selectedStatus = status;
                this.renderStatuses();
            });

            squaresDiv.appendChild(square);
        }

        this.statusView.appendChild(squaresDiv);
    }

    /**
     *
     * @param status
     */
    onStatusSelectClick(status) {
        this.currentStatus = status;

        this.sendMessage(MessageType.DRIVER_CHANGES_STATUS, {"status": status});

        this.renderStatuses();
    }

    renderOffers() {
        ClearE(this.offersView);

        this.offersView.appendChild(NewECT('h3', 'offers-title', 'Offers'));

        for (const offer of this.offers) {
            const parts = [
                NewECT('span', 'offer-client', offer.getPassengerId()),
                NewECT('span', 'offer-coords', `(${offer.lat}, ${offer.lng})`),
                NewECT('span', 'offer-status', offer.getStatus()),
            ];

            const actions = this.offerActions(offer);
            if (actions.hasChildNodes()) {
                parts.push(actions);
            }

            const row = NewEC('div', 'offer-row');
            parts.forEach((part, i) => {
                if (i > 0) {
                    row.appendChild(NewECT('span', 'offer-divider', '|'));
                }
                row.appendChild(part);
            });

            this.offersView.appendChild(row);
        }
    }

    offerActions(offer) {
        const cell = NewEC('span', 'offer-actions');

        const addButton = (text, onClick) => {
            const button = NewECT('button', 'offer-' + text, text);
            button.addEventListener('click', () => onClick(offer));
            cell.appendChild(button);
        };

        switch (offer.status) {
            case InvitationStatus.PENDING:
                addButton('accept', (o) => this.onOfferAcceptClick(o));
                addButton('reject', (o) => this.onOfferRejectClick(o));
                break;
            case InvitationStatus.ACCEPTED:
                addButton('cancel', (o) => this.onOfferCancelClick(o));
                break;
        }

        return cell;
    }

    onOfferAcceptClick(offer) {
        // TODO
    }

    onOfferRejectClick(offer) {
        // TODO
    }

    onOfferCancelClick(offer) {
        // TODO
    }

    initWs(isTls) {
        let protocol = 'ws://';

        if (isTls) {
            protocol = 'wss://';
        }

        if (!window["WebSocket"]) {
            this.driverView.appendChild(NewT('Your browser does not support WebSockets.'));
            return;
        }
        this.conn = new WebSocket(protocol + document.location.host + "/ws-driver");
        this.conn.onclose = (evt) => {
            this.driverView.appendChild(NewT('Websocket connection closed.'));
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
            case MessageType.SERVER_OFFERS_PASSENGER:
                this.handleServerOffersPassenger(envelope.data);
                break;
            case MessageType.SERVER_CANCELS_OFFER:
                this.handleServerCancelsOffer(envelope.data);
                break;
            case MessageType.SERVER_NOTIFIES_VOYAGE_STARTED:
                this.handleServerNotifiesVoyageStarted(envelope.data);
                break;
            case MessageType.SERVER_NOTIFIES_VOYAGE_FINISHED:
                this.handleServerNotifiesVoyageFinished(envelope.data);
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
                const lat = position.coords.latitude;
                const lon = position.coords.longitude;
                console.log('got coordinates', lat, lon);
                this.sendMessage(MessageType.CLIENT_RESPONDS_COORDINATES, {lat, lon});
            },
            (error) => {
                console.error('failed to get coordinates:', error.message);
                this.sendMessage(MessageType.CLIENT_ERROR, {error: error.message});
            },
            {enableHighAccuracy: true, maximumAge: 5000, timeout: 10000}
        );
    }

    sendMessage(type, data) {
        console.log('sending message', type, data);
        this.conn.send(JSON.stringify({
            id: crypto.randomUUID(),
            type: type,
            data: data,
        }));
    }

    handleServerRequestsClientInfo(data) {
        // TODO: real phone/vehicle_info once there's an input for them; empty for now.
        console.log('server requested client info, replying with empty values for now');
        this.sendMessage(MessageType.CLIENT_RESPONDS_INFO, {});
    }

    handleServerOffersPassenger(data) {
        console.log('server offered passenger, not implemented yet', data);
        // TODO: add data to the incoming-invites table
    }

    handleServerCancelsOffer(data) {
        console.log('server canceled offer, not implemented yet', data);
        // TODO: remove the offer from the incoming-invites table
    }

    handleServerNotifiesVoyageStarted(data) {
        console.log('server notified voyage started, not implemented yet', data);
        // TODO: switch UI to "trip in progress"
    }

    handleServerNotifiesVoyageFinished(data) {
        console.log('server notified voyage finished, not implemented yet', data);
        // TODO: switch UI back to available/idle
    }
}