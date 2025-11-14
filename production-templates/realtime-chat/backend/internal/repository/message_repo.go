package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"realtime-chat/internal/domain"
)

type MessageRepository interface {
	Create(message *domain.Message) error
	FindByID(id uint) (*domain.Message, error)
	FindByRoomID(roomID uint, limit, offset int) ([]*domain.Message, error)
	Update(message *domain.Message) error
	SoftDelete(messageID uint) error
	GetLastMessage(roomID uint) (*domain.Message, error)

	// Reaction operations
	AddReaction(reaction *domain.MessageReaction) error
	RemoveReaction(messageID, userID uint, emoji string) error
	GetReactions(messageID uint) ([]*domain.MessageReaction, error)

	// Read receipt operations
	MarkAsRead(messageID, userID uint) error
	GetReadReceipts(messageID uint) ([]*domain.ReadReceipt, error)
	GetLastReadMessage(roomID, userID uint) (*domain.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Reload message with sender
	return r.db.Preload("Sender").First(message, message.ID).Error
}

func (r *messageRepository) FindByID(id uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		Preload("ReadReceipts.User").
		First(&message, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("message not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}
	return &message, nil
}

func (r *messageRepository) FindByRoomID(roomID uint, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.Where("room_id = ? AND is_deleted = ?", roomID, false).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find messages for room: %w", err)
	}

	// Reverse order so oldest messages come first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) Update(message *domain.Message) error {
	if err := r.db.Save(message).Error; err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}

func (r *messageRepository) SoftDelete(messageID uint) error {
	if err := r.db.Model(&domain.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
			"content":    "[deleted]",
		}).Error; err != nil {
		return fmt.Errorf("failed to soft delete message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetLastMessage(roomID uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.Where("room_id = ? AND is_deleted = ?", roomID, false).
		Preload("Sender").
		Order("created_at DESC").
		First(&message).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No messages in room
		}
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}
	return &message, nil
}

// Reaction operations

func (r *messageRepository) AddReaction(reaction *domain.MessageReaction) error {
	// Check if reaction already exists
	var existing domain.MessageReaction
	err := r.db.Where("message_id = ? AND user_id = ? AND emoji = ?",
		reaction.MessageID, reaction.UserID, reaction.Emoji).
		First(&existing).Error

	if err == nil {
		// Reaction already exists
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing reaction: %w", err)
	}

	// Create new reaction
	if err := r.db.Create(reaction).Error; err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	// Reload with user
	return r.db.Preload("User").First(reaction, reaction.ID).Error
}

func (r *messageRepository) RemoveReaction(messageID, userID uint, emoji string) error {
	if err := r.db.Where("message_id = ? AND user_id = ? AND emoji = ?",
		messageID, userID, emoji).
		Delete(&domain.MessageReaction{}).Error; err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}
	return nil
}

func (r *messageRepository) GetReactions(messageID uint) ([]*domain.MessageReaction, error) {
	var reactions []*domain.MessageReaction
	err := r.db.Where("message_id = ?", messageID).
		Preload("User").
		Find(&reactions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}
	return reactions, nil
}

// Read receipt operations

func (r *messageRepository) MarkAsRead(messageID, userID uint) error {
	// Check if read receipt already exists
	var existing domain.ReadReceipt
	err := r.db.Where("message_id = ? AND user_id = ?", messageID, userID).
		First(&existing).Error

	if err == nil {
		// Already marked as read
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing read receipt: %w", err)
	}

	// Create read receipt
	receipt := &domain.ReadReceipt{
		MessageID: messageID,
		UserID:    userID,
	}

	if err := r.db.Create(receipt).Error; err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

func (r *messageRepository) GetReadReceipts(messageID uint) ([]*domain.ReadReceipt, error) {
	var receipts []*domain.ReadReceipt
	err := r.db.Where("message_id = ?", messageID).
		Preload("User").
		Order("read_at DESC").
		Find(&receipts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get read receipts: %w", err)
	}
	return receipts, nil
}

func (r *messageRepository) GetLastReadMessage(roomID, userID uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.
		Joins("JOIN read_receipts ON messages.id = read_receipts.message_id").
		Where("messages.room_id = ? AND read_receipts.user_id = ?", roomID, userID).
		Order("read_receipts.read_at DESC").
		First(&message).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last read message: %w", err)
	}
	return &message, nil
}
