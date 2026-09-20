package ws

import (
	"bytes"
	"log"
	"time"

	"darbelis.eu/taxi/internal/messages"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	clientId         string
	outgoingMessages chan ClientMessage
	conn             *websocket.Conn
	incomingMessages chan ClientMessage
	// will pass client id when closed
	closedNotifier chan string
}

func (c *Client) GetClientId() string {
	return c.clientId
}

func NewClient(clientId string, conn *websocket.Conn, incomingMessages chan ClientMessage, closedNotifier chan string) *Client {
	return &Client{
		clientId:         clientId,
		outgoingMessages: make(chan ClientMessage, 256),
		conn:             conn,
		incomingMessages: incomingMessages,
		closedNotifier:   closedNotifier,
	}
}

func (c *Client) Send(msg ClientMessage) {
	select {
	case c.outgoingMessages <- msg:
	default:
		c.closedNotifier <- c.clientId
	}
}

// ReadPump intended to run in a goroutine.
func (c *Client) readPump() {
	defer func() {
		c.closedNotifier <- c.clientId
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, newline, space, -1))

		messageId, message, err := messages.Decode(messageBytes)
		if err != nil {
			// better handling later
			log.Printf("error: %v", err)
			continue
		}

		clientMessage := ClientMessage{
			ClientId:  c.clientId,
			Message:   message,
			MessageId: messageId,
		}

		c.incomingMessages <- clientMessage
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case clientMessage, ok := <-c.outgoingMessages:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			messageBytes, err := messages.Encode(clientMessage.MessageId, clientMessage.Message)

			if err != nil {
				log.Printf("message encode error: %v", err)
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error initializing websocket writer: %v", err)

				return
			}
			w.Write(messageBytes)

			if err := w.Close(); err != nil {
				log.Printf("error closing websocket writer: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
