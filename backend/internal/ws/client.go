package ws

import (
	"backend/internal/model"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

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
		case "answer":
			var answer AnswerPayload
			err := json.Unmarshal(socketMsg.Data, &answer)
			if err != nil {
				log.Println("invalid answer payload:", err)
				continue
			}

			roomCode := c.roomCode
			playerID := c.playerID

			questionKey := fmt.Sprintf("questions:%s", roomCode)
			data, err := store.Client.Get(store.Ctx, questionKey).Result()
			if err != nil {
				log.Println("failed to fetch questions", err)
				return
			}

			var questions []model.Question
			err = json.Unmarshal([]byte(data), &questions)
			if err != nil {
				log.Println("failed to parse questions:", err)
				return
			}
			var question model.Question
			found := false
			for _, q := range questions {
				if q.ID == answer.QuestionID {
					question = q
					found = true
					break
				}
			}
			if !found {
				log.Println("question not found:", answer.QuestionID)
				return
			}

			scoreKey := fmt.Sprintf("score:%s:%s", roomCode, playerID)
			currentScore := 0
			rawScore, err := store.Client.Get(store.Ctx, scoreKey).Result()
			if err == nil {
				currentScore, _ = strconv.Atoi(rawScore)
			}

			if answer.Selected == question.CorrectAnswer {
				currentScore += 1000
			}

			store.Client.Set(store.Ctx, scoreKey, currentScore, 60*time.Minute)

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
