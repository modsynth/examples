package websocket

import (
	"log"
	"sync"
)

// Hub maintains active WebSocket connections and broadcasts messages
type Hub struct {
	// Registered clients organized by room ID
	rooms map[uint]map[*Client]bool

	// Inbound messages from clients
	broadcast chan *Message

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.RoomID] == nil {
		h.rooms[client.RoomID] = make(map[*Client]bool)
	}

	h.rooms[client.RoomID][client] = true
	log.Printf("Client registered: UserID=%d, RoomID=%d, Total in room=%d",
		client.UserID, client.RoomID, len(h.rooms[client.RoomID]))
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.RoomID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			// Remove room if no clients left
			if len(clients) == 0 {
				delete(h.rooms, client.RoomID)
			}

			log.Printf("Client unregistered: UserID=%d, RoomID=%d, Remaining in room=%d",
				client.UserID, client.RoomID, len(clients))
		}
	}
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[message.RoomID]
	if !ok {
		return
	}

	for client := range clients {
		// Don't send typing indicators back to the sender
		if message.Type == MessageTypeTyping && client.UserID == message.UserID {
			continue
		}

		select {
		case client.send <- message:
		default:
			// Client's send channel is full, close it
			close(client.send)
			delete(clients, client)
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// GetOnlineUsers returns a list of online user IDs in a room
func (h *Hub) GetOnlineUsers(roomID uint) []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make(map[uint]bool)
	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			userIDs[client.UserID] = true
		}
	}

	result := make([]uint, 0, len(userIDs))
	for userID := range userIDs {
		result = append(result, userID)
	}

	return result
}

// GetRoomCount returns the number of active rooms
func (h *Hub) GetRoomCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

// GetClientCount returns the total number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, clients := range h.rooms {
		count += len(clients)
	}
	return count
}
