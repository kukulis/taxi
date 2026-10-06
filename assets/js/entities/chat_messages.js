// Mirrors internal/messages/message_chat_payloads.go — keep the two in sync.
import {MessageType} from "../constants.js";

export class ChatMessage {
    constructor(id = '', senderId = '', receiverId = '', text = '', sentTime = '', receivedTime = '') {
        this.id = id;
        this.sender_id = senderId;
        this.receiver_id = receiverId;
        this.text = text;
        this.sent_time = sentTime;
        this.received_time = receivedTime;
    }

    fromObject(obj) {
        this.id = obj.id;
        this.sender_id = obj.sender_id;
        this.receiver_id = obj.receiver_id;
        this.text = obj.text;
        this.sent_time = obj.sent_time;
        this.received_time = obj.received_time;
        return this;
    }

    getMessageType() {
        return MessageType.CHAT_MESSAGE;
    }
}

export class ChatLog {
    constructor(id = '', senderId = '', receiverId = '', log = [], startedAt = '') {
        this.id = id;
        this.sender_id = senderId;
        this.receiver_id = receiverId;
        this.log = log;
        this.started_at = startedAt;
    }

    fromObject(obj) {
        this.id = obj.id;
        this.sender_id = obj.sender_id;
        this.receiver_id = obj.receiver_id;
        this.log = (obj.log || []).map((m) => new ChatMessage().fromObject(m));
        this.started_at = obj.started_at;
        return this;
    }

    getMessageType() {
        return MessageType.CHAT_LOG;
    }
}

export class ChatInfo {
    constructor(id = '', senderId = '', receiverId = '', startedAt = '') {
        this.id = id;
        this.sender_id = senderId;
        this.receiver_id = receiverId;
        this.started_at = startedAt;
    }

    fromObject(obj) {
        this.id = obj.id;
        this.sender_id = obj.sender_id;
        this.receiver_id = obj.receiver_id;
        this.started_at = obj.started_at;
        return this;
    }

    getMessageType() {
        return MessageType.CHAT_INFO;
    }
}

export class ChatInfos {
    constructor(list = []) {
        this.list = list;
    }

    fromObject(obj) {
        this.list = (obj.list || []).map((i) => new ChatInfo().fromObject(i));
        return this;
    }

    getMessageType() {
        return MessageType.CHAT_INFOS;
    }
}

export class ClientRequestsChatsInfos {
    fromObject(obj) {
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_REQUESTS_CHATS_INFOS;
    }
}

export class ClientRequestsChatLog {
    constructor(senderId = '', receiverId = '') {
        this.sender_id = senderId;
        this.receiver_id = receiverId;
    }

    fromObject(obj) {
        this.sender_id = obj.sender_id;
        this.receiver_id = obj.receiver_id;
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_REQUESTS_CHAT_LOG;
    }
}

export class TalkerInfo {
    constructor(clientId = '', clientType = '', name = '') {
        this.client_id = clientId;
        this.client_type = clientType;
        this.name = name;
    }

    fromObject(obj) {
        this.client_id = obj.client_id;
        this.client_type = obj.client_type;
        this.name = obj.name;
        return this;
    }

    getMessageType() {
        return MessageType.TALKER_INFO;
    }
}

export class TalkersInfos {
    constructor(list = []) {
        this.list = list;
    }

    fromObject(obj) {
        this.list = (obj.list || []).map((t) => new TalkerInfo().fromObject(t));
        return this;
    }

    getMessageType() {
        return MessageType.TALKERS_INFOS;
    }
}

export class ClientRequestsTalkersInfos {
    fromObject(obj) {
        return this;
    }

    getMessageType() {
        return MessageType.CLIENT_REQUESTS_TALKERS_INFOS;
    }
}
