import {NewEC, NewT} from "./util.js";
import {MessageType} from "./constants.js";

export class PassengerComponent {
    // main view
    passengerView = null;
    // websocket connection
    conn = null;

    async render() {
        this.passengerView = NewEC('div', 'passenger-component');

        this.passengerView.appendChild(NewT('TODO'));

        return this.passengerView;
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
