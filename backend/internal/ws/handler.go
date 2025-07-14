package ws

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid websocket path", http.StatusBadRequest)
		return
	}

	roomCode := parts[2]
	playerID := parts[3]
	log.Printf("Incoming WS: /ws/%s/%s\n", roomCode, playerID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		roomCode: roomCode,
		playerID: playerID,
	}
	GlobalHub.register <- client

	go client.writePump()
	go client.readPump()

}
