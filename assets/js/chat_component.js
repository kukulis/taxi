import {AppChildren, ClearE, NewE, NewEC, NewECT} from "./util.js";
import {Conversation} from "./entities/conversation.js";
import {ChatMessageDto} from "./entities/chat_message_dto.js";
import {Talker} from "./entities/talker.js";
import {TalkersList} from "./entities/talkers_list.js";
import {ApplicationEvent, MessageType} from "./constants.js";
import {ChatMessage, ChatLog, ChatInfo, ChatInfos, TalkerInfo, TalkersInfos} from "./entities/chat_messages.js";

export class ChatComponent {

    /**
     * @type {HTMLElement}
     */
    view = null;

    /**
     * @type {HTMLElement}
     */
    showButton = null;

    /**
     * @type {HTMLElement}
     */
    hideButton = null;

    /**
     * @type {Conversation}
     */
    conversation = null;

    /**
     * @type {HTMLElement}
     */
    messagesBox = null;

    /**
     * @type {TalkersList}
     */
    talkersList = null;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     *
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     *
     * @type {string}
     */
    clientId = null;

    /**
     *
     * @type {string}
     */
    clientType = null;

    /**
     * @type {Map<string, int>}
     */
    incomingChatNotices = new Map();

    /**
     * @type {HTMLElement}
     */
    incomingChatNoticesBox = null;

    /**
     * @type {State}
     */
    state = null;

    /**
     * @param {string} clientId
     * @param {string} clientType
     * @param {Dispatcher} dispatcher
     * @param {ApiClient} apiClient
     * @param {State} state
     */
    constructor(clientId,
                clientType,
                dispatcher,
                apiClient ,
                state, )
    {
        this.clientId = clientId;
        this.clientType = clientType;
        this.dispatcher = dispatcher;
        this.apiClient = apiClient;
        this.state = state;
    }

    render() {
        this.view = NewEC('div', 'messages-container');
        this.view.style.display = 'none';

        this.showButton = NewECT('button', 'chat-show-button', 'Chat');
        this.showButton.addEventListener('click', () => this.onShowButtonClick());

        this.hideButton = NewECT('button', 'chat-hide-button', 'Close chat');
        this.hideButton.addEventListener('click', () => this.onHideButtonClick());

        this.renderView();

        const root = NewEC('div', 'chat-component');
        return AppChildren(root, [this.showButton, this.view]);
    }

    onShowButtonClick() {
        this.view.style.display = '';
        this.showButton.style.display = 'none';
    }

    onHideButtonClick() {
        this.view.style.display = 'none';
        this.showButton.style.display = '';
    }

    renderView() {
        ClearE(this.view);
        this.view.appendChild(this.hideButton);
        this.messagesBox = null;

        this.view.appendChild(this.renderIncomingChatNotices());

        if (this.talkersList != null ) {
            const box = NewEC('div', 'chat-talkers-box');
            box.appendChild(NewECT('h3', 'chat-talkers-title', 'Talkers'));

            if (this.talkersList.talkers.length === 0) {
                box.appendChild(NewECT('span', 'chat-talkers-empty', 'empty'));
            }

            for (const talker of this.talkersList.talkers) {
                const row = NewEC('div', 'chat-talker-row');
                row.appendChild(NewECT('span', 'chat-talker-name', talker.name));

                const chatButton = NewECT('button', 'chat-talker-chat-button', 'Chat');
                chatButton.addEventListener('click', () => this.onTalkerChatButtonClick(talker));
                row.appendChild(chatButton);

                box.appendChild(row);
            }

            const closeButton = NewECT('button', 'chat-talkers-close-button', 'Close');
            closeButton.addEventListener('click', () => this.onTalkersCloseButtonClick());
            box.appendChild(closeButton);

            this.view.appendChild(box);
        }

        if ( this.conversation != null ) {
            const box = NewEC('div', 'chat-conversation-box');
            box.appendChild(NewECT('h3', 'chat-conversation-title', this.conversation.talkerName));

            const messagesBox = NewEC('div', 'chat-messages-box');
            for (const message of this.conversation.messages) {
                messagesBox.appendChild(this.createMessageRow(message));
            }
            box.appendChild(messagesBox);
            this.messagesBox = messagesBox;

            const inputRow = NewEC('div', 'chat-input-row');
            const input = NewE('input');
            input.classList.add('chat-input');
            input.type = 'text';
            inputRow.appendChild(input);

            const sendButton = NewECT('button', 'chat-send-button', 'Send');
            sendButton.addEventListener('click', () => this.onSendButtonClick(input));
            inputRow.appendChild(sendButton);
            box.appendChild(inputRow);

            const closeButton = NewECT('button', 'chat-close-button', 'Close chat');
            closeButton.addEventListener('click', () => this.onCloseButtonClick());
            box.appendChild(closeButton);

            this.view.appendChild(box);
        }

        if ( this.talkersList == null && this.conversation == null ) {
            const showTalkersButton = NewECT('button', 'chat-show-talkers-button', 'Show talkers');
            showTalkersButton.addEventListener('click', () => this.onShowTalkersButtonClick());
            this.view.appendChild(showTalkersButton);
        }
    }

    renderIncomingChatNotices() {
        this.incomingChatNoticesBox = NewEC('div', 'chat-incoming-notices-box');
        this.refreshIncomingChatNotices();

        return this.incomingChatNoticesBox;
    }

    refreshIncomingChatNotices() {
        const box = this.incomingChatNoticesBox;
        ClearE(box);
        box.appendChild(NewECT('h3', 'chat-incoming-notices-title', 'Incoming Chat Notices'));

        for (const [senderId, count] of this.incomingChatNotices) {
            const button = NewECT('button', 'chat-incoming-notice-button', `${senderId} (${count})`);
            button.addEventListener('click', () => {
                // fake talker, just to reuse onTalkerChatButtonClick - we will refactor that later
                const talker = new Talker(senderId, 'TODO');
                this.onTalkerChatButtonClick(talker);
            });
            box.appendChild(button);
        }
    }

    /**
     *
     * @param {ChatMessageDto} message
     * @returns {HTMLElement}
     */
    createMessageRow(message) {
        let senderName = this.state.name;
        if ( this.constructor != null && message.senderId == this.conversation.talkerId ) {
            senderName = this.conversation.talkerName
        }
        const row = NewEC('div', 'chat-message-row');
        row.appendChild(NewECT('span', 'chat-message-sender', senderName  ))
        row.appendChild(NewECT('span', 'chat-message-text', message.text));
        return row;
    }

    async onTalkerChatButtonClick(talker) {
        const conversation = new Conversation();
        conversation.talkerId = talker.id;
        conversation.talkerName = talker.name;

        this.conversation = conversation;
        this.talkersList = null;
        this.renderView();

        const chatLog = await this.apiClient.getChatLog(this.clientId, talker.id);

        // the user may have closed/switched the conversation while this was in flight
        if (chatLog == null || this.conversation !== conversation) {
            return;
        }

        conversation.messages = chatLog.log.map((message) => {
            const dto = new ChatMessageDto();
            dto.time = new Date(message.sent_time);
            dto.text = message.text;
            dto.senderId = message.sender_id;
            dto.receiverId = message.receiver_id;

            return dto;
        });
        this.renderView();
    }

    onSendButtonClick(input) {
        const text = input.value.trim();
        if (text === '') {
            return;
        }

        const sentTime = new Date();

        this.sendMessage(MessageType.CHAT_MESSAGE, {
            id: crypto.randomUUID(),
            sender_id: this.clientId,
            receiver_id: this.conversation.talkerId,
            text: text,
            sent_time: sentTime.toISOString(),
        });

        const message = new ChatMessageDto();
        message.time = sentTime;
        message.text = text;
        this.conversation.messages.push(message);

        if (this.messagesBox != null) {
            this.messagesBox.appendChild(this.createMessageRow(message));
        }

        input.value = '';
    }

    onCloseButtonClick() {
        this.conversation = null;
        this.renderView();
    }

    onTalkersCloseButtonClick() {
        this.talkersList = null;
        this.renderView();
    }

    async onShowTalkersButtonClick() {
        const talkerInfos = await this.apiClient.getTalkers();

        const talkersList = new TalkersList();
        talkersList.talkers = talkerInfos
            .filter((t) => t.client_id !== this.clientId)
            .map((t) => new Talker(t.client_id, t.name));

        this.talkersList = talkersList;
        this.renderView();
    }

    sendMessage(type, data) {
        this.dispatcher.dispatch(ApplicationEvent.SEND_WS_MESSAGE, [this.clientType, type, data]);
    }


    handleChatMessage(data) {
        const message = new ChatMessage().fromObject(data);

        if (this.conversation == null || this.conversation.talkerId !== message.sender_id) {
            const count = this.incomingChatNotices.get(message.sender_id) || 0;
            this.incomingChatNotices.set(message.sender_id, count + 1);
            if (this.incomingChatNoticesBox != null) {
                this.refreshIncomingChatNotices();
            }
            return;
        }

        const dto = new ChatMessageDto();
        dto.time = new Date(message.sent_time);
        dto.text = message.text;
        this.conversation.messages.push(dto);

        if (this.messagesBox != null) {
            this.messagesBox.appendChild(this.createMessageRow(dto));
        }
    }

    handleChatLog(data) {
        const chatLog = new ChatLog().fromObject(data);
        // TODO build this.conversation from chatLog, render
    }

    handleChatInfo(data) {
        const chatInfo = new ChatInfo().fromObject(data);
        // TODO add/update an entry in this.talkersList, render
    }

    handleChatInfos(data) {
        const chatInfos = new ChatInfos().fromObject(data);
        // TODO replace this.talkersList from chatInfos.list, render
    }

    handleTalkerInfo(data) {
        const talkerInfo = new TalkerInfo().fromObject(data);
        // TODO add/update an entry in this.talkersList, render
    }

    handleTalkersInfos(data) {
        const talkersInfos = new TalkersInfos().fromObject(data);
        // TODO replace this.talkersList from talkersInfos.list, render
    }
}