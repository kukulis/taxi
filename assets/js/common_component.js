import {ClearE, NewEC, NewECT} from "./util.js";
import {ApplicationEvent, MessageType} from "./constants.js";
import {ClientName, ClientRespondsInfo, ServerSendsClientAddress} from "./entities/messages.js";

// Location/address handling shared by DriverComponent and PassengerComponent.
// Deliberately has no reference to either of them; a .gohtml page wires them
// together through the shared Dispatcher, WsService and State instances.
export class CommonComponent {

    /**
     * @type {string}
     */
    clientType = null;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {State}
     */
    state = null;

    /**
     * @type {HTMLElement}
     */
    view = null;

    /**
     * @type {boolean}
     */
    editingName = false;

    /**
     * @param {string} clientType
     * @param {Dispatcher} dispatcher
     * @param {State} state
     */
    constructor(clientType, dispatcher, state) {
        this.clientType = clientType;
        this.dispatcher = dispatcher;
        this.state = state;
    }

    render() {
        this.view = NewEC('div', 'client-info');
        this.renderView();

        return this.view;
    }

    renderView() {
        ClearE(this.view);

        const {lat, lon} = this.state.getCoordinates();
        if (lat !== null && lon !== null) {
            this.view.appendChild(NewECT('div', 'client-location', `(${lat}, ${lon})`));
        } else {
            this.view.appendChild(NewECT('div', 'client-location', '...'));
        }

        this.view.appendChild(NewECT('div', 'client-address', this.state.getAddress()));

        if (this.editingName) {
            this.view.appendChild(NewECT('label', 'name-label', 'name'));
            let nameInput = NewEC('input', 'name-input');
            nameInput.value = this.state.name;
            nameInput.type = 'text';
            this.view.appendChild(nameInput);

            const saveName = () => {
                this.state.name = nameInput.value;

                this.sendMessage(MessageType.CLIENT_NAME, new ClientName(
                    this.state.name
                ));

                this.editingName = false;
                this.renderView();
            };

            nameInput.addEventListener('keydown', (e) => {
                if (e.key === 'Enter') {
                    saveName();
                } else if (e.key === 'Escape') {
                    this.editingName = false;
                    this.renderView();
                }
            });

            let nameButton = NewECT('button', 'set-name-button', 'save');
            this.view.appendChild(nameButton);

            nameButton.addEventListener('click', saveName);
        } else {
            this.view.appendChild(NewECT('label', 'name-label', 'name'));
            let nameSpan = NewECT('span', 'name-display', this.state.name);
            nameSpan.addEventListener('click', () => {
                this.editingName = true;
                this.renderView();
            });
            this.view.appendChild(nameSpan);
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
                this.state.setCoordinates(lat, lon);
                this.sendMessage(MessageType.CLIENT_RESPONDS_COORDINATES, {lat: lat, lon: lon});
                this.renderView();
            },
            (error) => {
                console.error('failed to get coordinates:', error.message);
                this.sendMessage(MessageType.CLIENT_ERROR, {error: error.message});
            },
            {enableHighAccuracy: true, maximumAge: 5000, timeout: 10000}
        );
    }

    handleServerRequestsClientInfo(data) {
        // TODO: real phone/vehicle_info once there's an input for them; empty for now.
        console.log('server requested client info, replying with empty values for now');
        this.sendMessage(MessageType.CLIENT_RESPONDS_INFO, {});
    }

    handleServerSendsClientAddress(data) {
        const message = new ServerSendsClientAddress().fromObject(data);
        this.state.setAddress(message.address);
        this.renderView();
    }

    handleClientName(data) {
        const message = new ClientName().fromObject(data);
        this.state.name = message.name;
        this.renderView();
    }

    sendMessage(type, data) {
        this.dispatcher.dispatch(ApplicationEvent.SEND_WS_MESSAGE, [this.clientType, type, data]);
    }
}
