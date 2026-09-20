// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Register requests from the clients.
	registerChannel chan *Client

	// Unregister requests from clients.
	unregisterChannel chan string

	incomingMessagesChannel chan ClientMessage
	outgoingMessagesChannel chan ClientMessage

	messageHandler MessageHandler
}

func NewHub(handler MessageHandler) *Hub {
	return &Hub{
		registerChannel:         make(chan *Client),
		unregisterChannel:       make(chan string, 256),
		clients:                 make(map[string]*Client),
		incomingMessagesChannel: make(chan ClientMessage, 256),
		outgoingMessagesChannel: make(chan ClientMessage, 256),
		messageHandler:          handler,
	}
}

// Run this makes registering, unregistering and message sending actions never happen in parallel
// need to find a way to get a list of clients from outside ... or do we?
// Can we mitigate a client's list in to an outside structure by using incoming messages from clients
// or adding hooks for various events?

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.registerChannel:
			h.clients[c.GetClientId()] = c
		case clientId := <-h.unregisterChannel:
			client, ok := h.clients[clientId]

			if !ok {
				log.Println("Client not found to close", clientId)
				break
			}
			delete(h.clients, clientId)

			// consider whether to use a wrapper Close function
			client.conn.Close()

		case clientMessage := <-h.outgoingMessagesChannel:
			client, ok := h.clients[clientMessage.ClientId]
			if !ok {
				log.Println("Client not found to send message to", clientMessage.ClientId)
				break
			}

			sentOk := client.send(clientMessage)
			if !sentOk {
				delete(h.clients, clientMessage.ClientId)
			}

		case clientMessage := <-h.incomingMessagesChannel:
			h.messageHandler.Handle(clientMessage)
		}
	}
}

// RegisterWebSocketClient handles websocket requests from the peer.
func (h *Hub) RegisterWebSocketClient(clientId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := NewClient(clientId, conn, h.incomingMessagesChannel, h.unregisterChannel)
	h.registerChannel <- c

	go c.writePump()
	go c.readPump()
}

func (h *Hub) SendMessage(message ClientMessage) {
	h.outgoingMessagesChannel <- message
}
