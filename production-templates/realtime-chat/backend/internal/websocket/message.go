package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Message events
	MessageTypeNewMessage    MessageType = "NEW_MESSAGE"
	MessageTypeMessageEdited MessageType = "MESSAGE_EDITED"
	MessageTypeMessageDeleted MessageType = "MESSAGE_DELETED"

	// Reaction events
	MessageTypeReactionAdded   MessageType = "REACTION_ADDED"
	MessageTypeReactionRemoved MessageType = "REACTION_REMOVED"

	// Typing indicator
	MessageTypeTyping MessageType = "TYPING"

	// Read receipts
	MessageTypeMessageRead MessageType = "MESSAGE_READ"

	// Room events
	MessageTypeUserJoined MessageType = "USER_JOINED"
	MessageTypeUserLeft   MessageType = "USER_LEFT"
	MessageTypeRoomUpdated MessageType = "ROOM_UPDATED"

	// User status
	MessageTypeUserStatusChanged MessageType = "USER_STATUS_CHANGED"

	// System messages
	MessageTypePing MessageType = "PING"
	MessageTypePong MessageType = "PONG"
	MessageTypeError MessageType = "ERROR"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	RoomID    uint        `json:"room_id"`
	UserID    uint        `json:"user_id"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewMessage creates a new WebSocket message
func NewMessage(msgType MessageType, roomID, userID uint, data interface{}) *Message {
	return &Message{
		Type:      msgType,
		RoomID:    roomID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ErrorMessage creates an error message
func ErrorMessage(roomID, userID uint, errorMsg string) *Message {
	return &Message{
		Type:      MessageTypeError,
		RoomID:    roomID,
		UserID:    userID,
		Data:      map[string]string{"error": errorMsg},
		Timestamp: time.Now(),
	}
}
