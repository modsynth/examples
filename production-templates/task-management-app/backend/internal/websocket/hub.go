package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type MessageType string

const (
	TypeTaskCreated MessageType = "TASK_CREATED"
	TypeTaskUpdated MessageType = "TASK_UPDATED"
	TypeTaskDeleted MessageType = "TASK_DELETED"
	TypeTaskMoved   MessageType = "TASK_MOVED"
	TypeCommentAdded MessageType = "COMMENT_ADDED"
	TypeUserJoined  MessageType = "USER_JOINED"
	TypeUserLeft    MessageType = "USER_LEFT"
)

type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	ProjectID uint        `json:"project_id"`
	UserID    uint        `json:"user_id"`
}

type Hub struct {
	// Project ID -> map of client connections
	projects   map[uint]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		projects:   make(map[uint]map[*Client]bool),
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

	if h.projects[client.ProjectID] == nil {
		h.projects[client.ProjectID] = make(map[*Client]bool)
	}
	h.projects[client.ProjectID][client] = true

	log.Printf("Client registered for project %d, total clients: %d",
		client.ProjectID, len(h.projects[client.ProjectID]))

	// Notify others that user joined
	go h.Broadcast(&Message{
		Type:      TypeUserJoined,
		ProjectID: client.ProjectID,
		UserID:    client.UserID,
		Payload: map[string]interface{}{
			"user_id": client.UserID,
		},
	})
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.projects[client.ProjectID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			if len(clients) == 0 {
				delete(h.projects, client.ProjectID)
			}

			log.Printf("Client unregistered from project %d, remaining: %d",
				client.ProjectID, len(clients))

			// Notify others that user left
			go h.Broadcast(&Message{
				Type:      TypeUserLeft,
				ProjectID: client.ProjectID,
				UserID:    client.UserID,
				Payload: map[string]interface{}{
					"user_id": client.UserID,
				},
			})
		}
	}
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.projects[message.ProjectID]
	if !ok {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		// Don't send message back to sender
		if client.UserID == message.UserID {
			continue
		}

		select {
		case client.send <- data:
		default:
			// Client's send channel is full, remove it
			close(client.send)
			delete(clients, client)
		}
	}
}

func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

func (h *Hub) GetOnlineUsers(projectID uint) []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.projects[projectID]
	if !ok {
		return []uint{}
	}

	userIDs := make(map[uint]bool)
	for client := range clients {
		userIDs[client.UserID] = true
	}

	result := make([]uint, 0, len(userIDs))
	for userID := range userIDs {
		result = append(result, userID)
	}

	return result
}
