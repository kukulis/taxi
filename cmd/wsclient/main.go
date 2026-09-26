package main

import (
	"encoding/json"
	"flag"
	"log"

	"darbelis.eu/taxi/internal/message_common"
	"darbelis.eu/taxi/internal/messages"
	"darbelis.eu/taxi/pkg/util"
	"github.com/gorilla/websocket"
)

func main() {
	url := flag.String("url", "ws://localhost:8880/ws-driver", "WebSocket endpoint to connect to")
	msgType := flag.String("type", "", "message type, e.g. driver_accepts_offer (see internal/messages/message_types.go)")
	msgId := flag.String("id", "", "message id; a random one is generated when omitted")
	data := flag.String("data", "{}", "payload JSON, without the envelope, e.g. {\"passenger_id\":\"p1\"}")
	flag.Parse()

	if *msgType == "" {
		log.Fatal("-type is required")
	}
	if !json.Valid([]byte(*data)) {
		log.Fatalf("-data is not valid JSON: %s", *data)
	}

	id := *msgId
	if id == "" {
		var err error
		id, err = util.RandomHash(8)
		if err != nil {
			log.Fatalf("failed to generate a message id: %v", err)
		}
	}

	envelope := messages.Envelope{
		Id:   id,
		Type: message_common.MessageType(*msgType),
		Data: json.RawMessage(*data),
	}

	raw, err := json.Marshal(envelope)
	if err != nil {
		log.Fatalf("failed to marshal envelope: %v", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(*url, nil)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", *url, err)
	}
	defer conn.Close()

	log.Printf("connected to %s", *url)
	log.Printf("sending: %s", raw)

	if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
		log.Fatalf("failed to write message: %v", err)
	}

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("connection closed: %v", err)
			return
		}
		log.Printf("received: %s", messageBytes)
	}
}
