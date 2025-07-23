package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
// Each connected player gets their own Client instance.
//
// The client is associated with a specific room (by roomCode)
// and communicates with the Hub via send and receive channels.
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	roomCode string
	playerID string
}

type SocketMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type AnswerPayload struct {
	QuestionID string `json:"questionId"`
	Selected   string `json:"selected"`
}

// readPump listens for incoming WebSocket messages from the client.
// It should run as a goroutine per connection.
// When the client disconnects or an error occurs, it cleans up the connection.
func (c *Client) readPump() {
	defer c.conn.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var socketMsg SocketMessage
		err = json.Unmarshal(msg, &socketMsg)
		if err != nil {
			log.Println("invalid socket message:", err)
			continue
		}

		switch socketMsg.Type {
		default:
			log.Println("unknown message type:", socketMsg.Type)
		}
	}
}

// writePump listens on the send channel and writes messages to the WebSocket connection.
// It should be started as a goroutine for each client.
// If sending fails or the connection is closed, it cleans up the connection.
func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}

}
