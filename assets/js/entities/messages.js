// Mirrors internal/messages/message_payloads.go — keep the two in sync.
// Field names are the Go json tags (snake_case), so an instance can be sent as-is
// and an incoming envelope.data can be loaded with new X().fromObject(data).
import {MessageType} from "../constants.js";

// === to driver or to passenger

export class ServerRequestsCoordinates {
    fromObject(obj) {
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_REQUESTS_COORDINATES;
    }
}

export class ServerRequestsClientInfo {
    fromObject(obj) {
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_REQUESTS_CLIENT_INFO;
    }
}

// === driver or passenger

export class ClientRespondsCoordinates {
    constructor(lat = 0, lon = 0) {
        this.lat = lat;
        this.lon = lon;
    }

    fromObject(obj) {
        this.lat = obj.lat;
        this.lon = obj.lon;
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_RESPONDS_COORDINATES;
    }
}

// vehicle_info is only meaningful for a driver; ip is optional (server keeps the one on record when empty).
export class ClientRespondsInfo {
    constructor(phone = '', vehicleInfo = '', ip = '') {
        this.phone = phone;
        this.vehicle_info = vehicleInfo;
        this.ip = ip;
    }

    fromObject(obj) {
        this.phone = obj.phone;
        this.vehicle_info = obj.vehicle_info;
        this.ip = obj.ip;
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_RESPONDS_INFO;
    }
}

// sent when the client didn't understand a server request
export class ClientResponseError {
    constructor(error = '') {
        this.error = error;
    }

    fromObject(obj) {
        this.error = obj.error;
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_RESPONSE_ERROR;
    }
}

// sent when the client hit its own error, unrelated to a server request
export class ClientError {
    constructor(error = '') {
        this.error = error;
    }

    fromObject(obj) {
        this.error = obj.error;
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_ERROR;
    }
}

// --- driver-related messages

export class ServerOffersPassenger {
    constructor(passengerId = '', lat = 0, lon = 0) {
        this.passenger_id = passengerId;
        this.lat = lat;
        this.lon = lon;
    }

    fromObject(obj) {
        this.passenger_id = obj.passenger_id;
        this.lat = obj.lat;
        this.lon = obj.lon;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_OFFERS_PASSENGER;
    }
}

export class ServerCancelsOffer {
    constructor(passengerId = '') {
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_CANCELS_OFFER;
    }
}

export class DriverCancelsOffer {
    constructor(passengerId = '') {
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.DRIVER_CANCELS_OFFER;
    }
}

export class DriverAcceptsOffer {
    constructor(passengerId = '') {
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.DRIVER_ACCEPTS_OFFER;
    }
}

export class DriverRejectsOffer {
    constructor(passengerId = '') {
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.DRIVER_REJECTS_OFFER;
    }
}

export class DriverChangesStatus {
    // one of DriverStatus
    constructor(status = '') {
        this.status = status;
    }

    fromObject(obj) {
        this.status = obj.status;
        return this;
    }

    getMessageType() {
        return MessageType.DRIVER_CHANGES_STATUS;
    }
}

// --- passenger-related messages

export class PassengerInvitesDriver {
    constructor(driverId = '', lat = 0, lon = 0) {
        this.driver_id = driverId;
        this.lat = lat;
        this.lon = lon;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        this.lat = obj.lat;
        this.lon = obj.lon;
        return this;
    }

    getMessageType() {
        return MessageType.PASSENGER_INVITES_DRIVER;
    }
}

export class PassengerCancelsInvite {
    constructor(driverId = '') {
        this.driver_id = driverId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        return this;
    }

    getMessageType() {
        return MessageType.PASSENGER_CANCELS_INVITE;
    }
}

export class ServerNotifiesInviteAccepted {
    constructor(driverId = '') {
        this.driver_id = driverId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_NOTIFIES_INVITE_ACCEPTED;
    }
}

export class ServerNotifiesInviteRejected {
    constructor(driverId = '') {
        this.driver_id = driverId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_NOTIFIES_INVITE_REJECTED;
    }
}

export class ServerNotifiesInviteCanceled {
    constructor(driverId = '') {
        this.driver_id = driverId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_NOTIFIES_INVITE_CANCELED;
    }
}

export class PassengerRequestDriverCoords {
    constructor(driverId = '') {
        this.driver_id = driverId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        return this;
    }

    getMessageType() {
        return MessageType.PASSENGER_REQUEST_DRIVER_COORDS;
    }
}

export class ServerRespondsDriverCoords {
    constructor(driverId = '', lat = 0, lon = 0) {
        this.driver_id = driverId;
        this.lat = lat;
        this.lon = lon;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        this.lat = obj.lat;
        this.lon = obj.lon;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_RESPONDS_DRIVER_COORDS;
    }
}

export class ServerNotifiesVoyageStarted {
    constructor(driverId = '', passengerId = '') {
        this.driver_id = driverId;
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_NOTIFIES_VOYAGE_STARTED;
    }
}

export class ServerNotifiesVoyageFinished {
    constructor(driverId = '', passengerId = '') {
        this.driver_id = driverId;
        this.passenger_id = passengerId;
    }

    fromObject(obj) {
        this.driver_id = obj.driver_id;
        this.passenger_id = obj.passenger_id;
        return this;
    }

    getMessageType() {
        return MessageType.SERVER_NOTIFIES_VOYAGE_FINISHED;
    }
}
