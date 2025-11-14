package domain

import "time"

type RoomType string

const (
	RoomTypeDirect RoomType = "direct"   // 1:1 private chat
	RoomTypeGroup  RoomType = "group"    // Group chat
	RoomTypePublic RoomType = "public"   // Public channel
)

type Room struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         RoomType      `json:"type" gorm:"not null;default:'group'"`
	AvatarURL    string        `json:"avatar_url"`
	CreatorID    uint          `json:"creator_id" gorm:"not null"`
	Creator      *User         `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Participants []Participant `json:"participants,omitempty" gorm:"foreignKey:RoomID"`
	LastMessage  *Message      `json:"last_message,omitempty" gorm:"-"` // Not stored in DB, loaded separately
	IsArchived   bool          `json:"is_archived" gorm:"not null;default:false"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Participant struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoomID       uint      `json:"room_id" gorm:"not null;uniqueIndex:idx_room_user"`
	UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_room_user"`
	User         *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role         string    `json:"role" gorm:"not null;default:'member'"` // admin, member
	IsMuted      bool      `json:"is_muted" gorm:"not null;default:false"`
	LastReadAt   time.Time `json:"last_read_at"`
	UnreadCount  int       `json:"unread_count" gorm:"-"` // Calculated field
	JoinedAt     time.Time `json:"joined_at"`
	LeftAt       *time.Time `json:"left_at"`
}

type CreateRoomRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Type        RoomType `json:"type" binding:"required"`
	UserIDs     []uint   `json:"user_ids"` // Initial participants
}

type UpdateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

type AddParticipantRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role"` // admin or member
}
