import {WsEvent} from "./constants.js";

export class WsService {

    /**
     * @type {?string}
     */
    uri = null;


    /**
     * @type {boolean}
     */
    isTls = false;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    connection = null;


    /**
     * @type {Map<string, function>}
     */
    handles = new Map();

    /**
     *
     * @param {string} uri
     * @param {boolean} isTls
     * @param {Dispatcher} dispatcher
     */
    constructor(uri, isTls, dispatcher) {
        this.uri = uri;
        this.isTls = isTls;
        this.dispatcher = dispatcher;
    }

    /**
     *
     * @param {string} messageType
     * @param {function} handle
     */
    addHandle(messageType, handle) {
        this.handles.set(messageType, handle)
    }

    startWsHandling() {
        let protocol = 'ws://';

        if (this.isTls) {
            protocol = 'wss://';
        }

        if (!window["WebSocket"]) {
            this.dispatcher.dispatch('error', {message: 'Your browser does not support WebSockets.'});
            // this.driverView.appendChild(NewT('Your browser does not support WebSockets.'));
            return;
        }

        this.connection = new WebSocket(protocol + document.location.host + "/ws-driver");
        // send our coordinates right away instead of waiting for the server to ask;
        // only once open, since send() throws while the socket is still connecting
        this.connection.onopen = () => {
            this.dispatcher.dispatch(WsEvent.ON_OPEN_WS, {})
            // this.handleServerRequestsCoordinates({});
        };
        this.connection.onclose = (evt) => {
            // this.driverView.appendChild(NewT('Websocket connection closed.'));
            this.dispatcher.dispatch(WsEvent.ON_CLOSE_WS, {})
        };
        this.connection.onmessage = (evt) => {
            let messages = evt.data.split('\n');
            for (let i = 0; i < messages.length; i++) {
                const envelope = JSON.parse(messages[i]);

                if (this.handles.has(envelope.type)) {
                    this.handles.get(envelope.type)(envelope.data)
                    continue;
                }
                this.dispatcher.dispatch('error', {message: 'WS message type ' + envelope.type + ' unimplemented'})
            }
        };
    }

    sendMessage(type, data) {
        this.connection.send(JSON.stringify({
            id: crypto.randomUUID(),
            type: type,
            data: data,
        }));
    }
}