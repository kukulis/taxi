export class ChatMessageDto {
    /**
     * @type {Date}
     */
    time = null;
    senderId = "";
    receiverId = "";
    text = "";

    senderName = "";
    receiverName = "";
    sentTime = "";
    receivedTime = "";


    /**
     *
     * @param {ChatMessage} obj
     * @returns {ChatMessageDto}
     */
    fromWsMessage(obj) {
        this.id = obj.id;
        this.senderId = obj.sender_id;
        this.receiverId = obj.receiver_id;
        this.text = obj.text;
        this.sentTime = obj.sent_time;
        this.receivedTime = obj.received_time;

        this.senderName = obj.sender_name;
        this.receiverName = obj.receiver_name;

        return this;
    }
}