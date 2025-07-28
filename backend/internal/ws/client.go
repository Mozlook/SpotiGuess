package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10 // co ~54s
	writeWait  = 10 * time.Second
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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in readPump: %v", r)
		}
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("unexpected close error: %v", err)
			} else {
				log.Printf("client disconnected: %v", err)
			}
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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in writePump: %v", r)
		}
		c.conn.Close()
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write error:", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("ping error:", err)
				return
			}
		}
	}
}
