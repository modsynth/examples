package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"realtime-chat/internal/domain"
)

type RoomRepository interface {
	Create(room *domain.Room) error
	FindByID(id uint) (*domain.Room, error)
	FindByUserID(userID uint) ([]*domain.Room, error)
	FindDirectRoom(user1ID, user2ID uint) (*domain.Room, error)
	Update(room *domain.Room) error
	Delete(id uint) error

	// Participant operations
	AddParticipant(participant *domain.Participant) error
	RemoveParticipant(roomID, userID uint) error
	FindParticipant(roomID, userID uint) (*domain.Participant, error)
	GetParticipants(roomID uint) ([]*domain.Participant, error)
	UpdateLastRead(roomID, userID uint) error
	GetUnreadCount(roomID, userID uint) (int64, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *domain.Room) error {
	if err := r.db.Create(room).Error; err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *roomRepository) FindByID(id uint) (*domain.Room, error) {
	var room domain.Room
	err := r.db.Preload("Creator").
		Preload("Participants.User").
		First(&room, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find room: %w", err)
	}
	return &room, nil
}

func (r *roomRepository) FindByUserID(userID uint) ([]*domain.Room, error) {
	var rooms []*domain.Room

	err := r.db.
		Joins("JOIN participants ON rooms.id = participants.room_id").
		Where("participants.user_id = ? AND participants.left_at IS NULL", userID).
		Preload("Creator").
		Preload("Participants.User").
		Order("rooms.updated_at DESC").
		Find(&rooms).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find rooms for user: %w", err)
	}
	return rooms, nil
}

func (r *roomRepository) FindDirectRoom(user1ID, user2ID uint) (*domain.Room, error) {
	var room domain.Room

	// Find direct room with exactly these two users
	err := r.db.
		Where("type = ?", domain.RoomTypeDirect).
		Where("id IN (?)",
			r.db.Table("participants").
				Select("room_id").
				Where("user_id IN (?, ?)", user1ID, user2ID).
				Group("room_id").
				Having("COUNT(DISTINCT user_id) = 2"),
		).
		Preload("Creator").
		Preload("Participants.User").
		First(&room).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not an error, just no direct room exists
		}
		return nil, fmt.Errorf("failed to find direct room: %w", err)
	}
	return &room, nil
}

func (r *roomRepository) Update(room *domain.Room) error {
	if err := r.db.Save(room).Error; err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	return nil
}

func (r *roomRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

// Participant operations

func (r *roomRepository) AddParticipant(participant *domain.Participant) error {
	if err := r.db.Create(participant).Error; err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}

func (r *roomRepository) RemoveParticipant(roomID, userID uint) error {
	// Soft delete by setting left_at
	if err := r.db.Model(&domain.Participant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("left_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	return nil
}

func (r *roomRepository) FindParticipant(roomID, userID uint) (*domain.Participant, error) {
	var participant domain.Participant
	err := r.db.Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).
		Preload("User").
		First(&participant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("participant not found")
		}
		return nil, fmt.Errorf("failed to find participant: %w", err)
	}
	return &participant, nil
}

func (r *roomRepository) GetParticipants(roomID uint) ([]*domain.Participant, error) {
	var participants []*domain.Participant
	err := r.db.Where("room_id = ? AND left_at IS NULL", roomID).
		Preload("User").
		Find(&participants).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	return participants, nil
}

func (r *roomRepository) UpdateLastRead(roomID, userID uint) error {
	if err := r.db.Model(&domain.Participant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("last_read_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to update last read: %w", err)
	}
	return nil
}

func (r *roomRepository) GetUnreadCount(roomID, userID uint) (int64, error) {
	var participant domain.Participant
	err := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&participant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to find participant: %w", err)
	}

	var count int64
	err = r.db.Model(&domain.Message{}).
		Where("room_id = ? AND created_at > ? AND sender_id != ?",
			roomID, participant.LastReadAt, userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return count, nil
}
