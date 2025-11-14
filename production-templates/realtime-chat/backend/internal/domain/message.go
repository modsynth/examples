package domain

import "time"

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeSystem MessageType = "system" // System messages (user joined, left, etc.)
)

type Message struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	RoomID          uint              `json:"room_id" gorm:"not null;index"`
	SenderID        uint              `json:"sender_id" gorm:"not null"`
	Sender          *User             `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Type            MessageType       `json:"type" gorm:"not null;default:'text'"`
	Content         string            `json:"content"`
	FileURL         string            `json:"file_url"`
	FileName        string            `json:"file_name"`
	FileSize        int64             `json:"file_size"`
	FileMimeType    string            `json:"file_mime_type"`
	ReplyToID       *uint             `json:"reply_to_id"`
	ReplyTo         *Message          `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
	IsEdited        bool              `json:"is_edited" gorm:"not null;default:false"`
	EditedAt        *time.Time        `json:"edited_at"`
	IsDeleted       bool              `json:"is_deleted" gorm:"not null;default:false"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	Reactions       []MessageReaction `json:"reactions,omitempty" gorm:"foreignKey:MessageID"`
	ReadReceipts    []ReadReceipt     `json:"read_receipts,omitempty" gorm:"foreignKey:MessageID"`
	CreatedAt       time.Time         `json:"created_at" gorm:"index"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type MessageReaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MessageID uint      `json:"message_id" gorm:"not null;uniqueIndex:idx_message_user_reaction"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_message_user_reaction"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Emoji     string    `json:"emoji" gorm:"not null;uniqueIndex:idx_message_user_reaction"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadReceipt struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MessageID uint      `json:"message_id" gorm:"not null;uniqueIndex:idx_message_user_read"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_message_user_read"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ReadAt    time.Time `json:"read_at"`
}

type SendMessageRequest struct {
	Content   string      `json:"content"`
	Type      MessageType `json:"type"`
	ReplyToID *uint       `json:"reply_to_id"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type AddReactionRequest struct {
	Emoji string `json:"emoji" binding:"required"`
}

type TypingIndicator struct {
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IsTyping  bool      `json:"is_typing"`
	Timestamp time.Time `json:"timestamp"`
}
