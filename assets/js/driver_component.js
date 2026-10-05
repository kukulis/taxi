import {AppChildren, ClearE, NewEC, NewECT, NewT} from "./util.js";
import {ApplicationEvent, ClientType, DriverStatus, InvitationStatus, MessageType, WsEvent} from "./constants.js";
import {Offer} from "./entities/offer.js";
import {
    DriverAcceptsOffer,
    DriverCancelsOffer,
    DriverRejectsOffer,
    ServerCancelsOffer,
    ServerNotifiesVoyageFinished,
    ServerNotifiesVoyageStarted,
    ServerOffersPassenger,
    ServerRefreshDriverOffers,
    ServerRefreshDriverStatus,
} from "./entities/messages.js";

export class DriverComponent {

    driverId = null;
    // main view
    driverView = null;

    /**
     * @type {HTMLElement}
     */
    infoView = null;
    // status squares container
    statusView = null;
    // offers table container
    offersView = null;
    // websocket connection
    // conn = null;

    currentStatus = DriverStatus.RESTING;
    // status square the user clicked on, but has not confirmed with 'select' yet
    selectedStatus = null;

    /**
     *
     * @type {Offer[]}
     */
    offers = [];

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    constructor(driverId, dispatcher) {
        this.driverId = driverId;
        this.dispatcher = dispatcher;
    }

    async render() {
        this.driverView = NewEC('div', 'driver-component');
        this.infoView = NewEC('div', 'driver-info');
        this.statusView = NewEC('div', 'driver-statuses');
        this.offersView = NewEC('div', 'driver-offers');

        this.renderInfo();
        this.renderStatuses();
        this.renderOffers();

        return AppChildren(this.driverView, [this.infoView, this.statusView, this.offersView]);
    }

    renderInfo() {
        ClearE(this.infoView);
        this.infoView.appendChild(NewECT('h3', 'driver-title', 'Driver'));
        this.infoView.appendChild(NewECT('div', 'driver-id-info', this.driverId));
    }

    renderStatuses() {
        ClearE(this.statusView);


        const squaresDiv = NewEC('div', 'status-squares-container');
        for (const status of [DriverStatus.IDLE, DriverStatus.WORKING, DriverStatus.RESTING]) {
            const square = NewEC('div', 'status-square');
            square.classList.add('status-' + status);
            square.appendChild(NewECT('div', 'status-label', status));

            if (status === this.selectedStatus) {
                square.classList.add('status-square-selected');
            }

            if (status === this.currentStatus) {
                square.appendChild(NewECT('div', 'status-check', '✓'));
            } else if (status === this.selectedStatus) {
                const button = NewECT('button', 'status-select', 'Go');
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

    renderOffers() {
        ClearE(this.offersView);

        this.offersView.appendChild(NewECT('h3', 'offers-title', 'Offers'));

        for (const offer of this.offers) {
            const parts = [
                NewECT('span', 'offer-client', offer.getPassengerId()),
                NewECT('span', 'offer-coords', `(${offer.lat}, ${offer.lon})`),
                this.statusBadge('offer-status', offer.getStatus()),
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
            case InvitationStatus.ACCEPTED: {
                // invisible stand-in keeps the 'accept' slot occupied, so 'cancel' appears where 'reject' was
                // instead of sliding under the driver's finger right after tapping 'accept'
                const placeholder = NewECT('button', 'offer-placeholder', 'accept');
                placeholder.disabled = true;
                placeholder.setAttribute('aria-hidden', 'true');
                cell.appendChild(placeholder);
                addButton('cancel', (o) => this.onOfferCancelClick(o));
                break;
            }
            case InvitationStatus.REJECTED:
            case InvitationStatus.CANCELED:
            case InvitationStatus.COMPLETED:
                addButton('hide', (o) => this.onOfferHideClick(o));
                break;
        }

        return cell;
    }

    /**
     * Removes a finished offer from the list; the server keeps it.
     * @param {Offer} offer
     */
    onOfferHideClick(offer) {
        this.offers = this.offers.filter((o) => o !== offer);
        this.renderOffers();
    }

    /**
     * @param {Offer} offer
     */
    onOfferAcceptClick(offer) {
        const message = new DriverAcceptsOffer(offer.passengerId)
        this.sendMessage(message.getMessageType(), message);
        // is this pointer?
        offer.setStatus(InvitationStatus.ACCEPTED);

        this.renderOffers();
    }

    onOfferRejectClick(offer) {
        const message = new DriverRejectsOffer(offer.passengerId)
        this.sendMessage(message.getMessageType(), message);

        this.offers = this.offers.filter((o) => o !== offer);

        this.renderOffers();
    }

    onOfferCancelClick(offer) {
        const message = new DriverCancelsOffer(offer.passengerId)
        this.sendMessage(message.getMessageType(), message);

        this.offers = this.offers.filter((o) => o !== offer);

        this.renderOffers();
    }

    sendMessage(type, data) {
        console.log('driver component dispatches send ws message', type, data);
        this.dispatcher.dispatch(ApplicationEvent.SEND_WS_MESSAGE, [ClientType.DRIVER, type, data])
    }

    handleServerOffersPassenger(data) {
        const message = new ServerOffersPassenger().fromObject(data);

        const offer = new Offer()
            .setLat(message.lat)
            .setLng(message.lon)
            .setPassengerId(message.passenger_id)
            .setStatus(InvitationStatus.PENDING)
        ;

        this.offers.push(offer)

        this.renderOffers();
    }

    handleServerCancelsOffer(data) {
        const message = new ServerCancelsOffer().fromObject(data);
        const offer = this.offers.findLast((o) => o.getPassengerId() === message.passenger_id);
        if (!offer) {
            console.warn('no offer found for passenger', message.passenger_id);
            return;
        }
        offer.setStatus(InvitationStatus.CANCELED);
        this.renderOffers();
    }

    handleServerNotifiesVoyageStarted(data) {
        console.log('voyage started', data);

        const message = new ServerNotifiesVoyageStarted().fromObject(data);

        // latest accepted offer: older offers from the same passenger may still be in the list
        const offer = this.offers.findLast((o) =>
            o.getPassengerId() === message.passenger_id && o.getStatus() === InvitationStatus.ACCEPTED);
        if (!offer) {
            console.warn('no accepted offer for passenger', message.passenger_id);
            return;
        }

        offer.setStatus(InvitationStatus.DRIVING);
        this.renderOffers();
    }

    handleServerNotifiesVoyageFinished(data) {
        console.log('voyage finished', data);
        const message = new ServerNotifiesVoyageFinished().fromObject(data);

        // latest driving offer: older offers from the same passenger may still be in the list
        const offer = this.offers.findLast((o) =>
            o.getPassengerId() === message.passenger_id && o.getStatus() === InvitationStatus.DRIVING);
        if (!offer) {
            console.warn('no driving offer for passenger', message.passenger_id);
            return;
        }

        offer.setStatus(InvitationStatus.COMPLETED);
        this.renderOffers();
    }

    handleServerRefreshStatus(data) {
        const message = new ServerRefreshDriverStatus().fromObject(data)
        this.currentStatus = message.status
        this.renderStatuses()

    }

    handleServerRefreshOffers(data) {
        const message = new ServerRefreshDriverOffers().fromObject(data)
        this.offers = message.offers
        this.renderOffers()
    }
}