import {NewEC, NewT} from "./util.js";
import {MessageType} from "./constants.js";

export class DriverComponent {
    // main view
    driverView = null;
    // websocket connection
    conn = null;

    async render() {
        this.driverView = NewEC('div', 'driver-component');

        this.driverView.appendChild(NewT('TODO'));

        return this.driverView;
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
        console.log('server requested client info, not implemented yet', data);
        // TODO: send back client_responds_info
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