package ws

// Hub manages all active WebSocket clients, grouped by roomCode.
// It handles client registration, unregistration and broadcasting messages to all clients in a room.
type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Broadcast  chan BroadcastMessage
}

type BroadcastMessage struct {
	RoomCode string
	Data     []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
	}
}

var GlobalHub = NewHub()

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			clients, ok := h.rooms[client.roomCode]
			if !ok {
				clients = make(map[*Client]bool)
				h.rooms[client.roomCode] = clients
			}
			clients[client] = true

		case client := <-h.unregister:
			if clients, ok := h.rooms[client.roomCode]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.rooms, client.roomCode)
					}
				}
			}

		case msg := <-h.Broadcast:
			if clients, ok := h.rooms[msg.RoomCode]; ok {
				for client := range clients {
					select {
					case client.send <- msg.Data:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}

		}
	}

}
