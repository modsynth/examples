package service

import (
	"errors"
	"fmt"
	"time"

	"realtime-chat/internal/domain"
	"realtime-chat/internal/repository"
	"realtime-chat/internal/websocket"
)

type MessageService interface {
	Send(roomID, senderID uint, req *domain.SendMessageRequest) (*domain.Message, error)
	GetByID(messageID, userID uint) (*domain.Message, error)
	GetRoomMessages(roomID, userID uint, limit, offset int) ([]*domain.Message, error)
	Update(messageID, userID uint, req *domain.UpdateMessageRequest) (*domain.Message, error)
	Delete(messageID, userID uint) error

	// Reactions
	AddReaction(messageID, userID uint, req *domain.AddReactionRequest) error
	RemoveReaction(messageID, userID uint, emoji string) error

	// Read receipts
	MarkAsRead(messageID, userID uint) error

	// Typing indicator
	SendTypingIndicator(roomID, userID uint, isTyping bool) error
}

type messageService struct {
	messageRepo repository.MessageRepository
	roomRepo    repository.RoomRepository
	userRepo    repository.UserRepository
	hub         *websocket.Hub
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	roomRepo repository.RoomRepository,
	userRepo repository.UserRepository,
	hub *websocket.Hub,
) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		hub:         hub,
	}
}

func (s *messageService) Send(roomID, senderID uint, req *domain.SendMessageRequest) (*domain.Message, error) {
	// Verify sender is participant
	participant, err := s.roomRepo.FindParticipant(roomID, senderID)
	if err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	// Check if muted
	if participant.IsMuted {
		return nil, errors.New("you are muted in this room")
	}

	// Validate message content
	if req.Content == "" && req.Type == domain.MessageTypeText {
		return nil, errors.New("message content is required")
	}

	// Create message
	message := &domain.Message{
		RoomID:    roomID,
		SenderID:  senderID,
		Type:      req.Type,
		Content:   req.Content,
		ReplyToID: req.ReplyToID,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Reload message with sender and reply-to
	message, err = s.messageRepo.FindByID(message.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload message: %w", err)
	}

	// Broadcast new message event
	s.broadcastMessageEvent(roomID, senderID, websocket.MessageTypeNewMessage, message)

	return message, nil
}

func (s *messageService) GetByID(messageID, userID uint) (*domain.Message, error) {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Verify user has access to this message's room
	if _, err := s.roomRepo.FindParticipant(message.RoomID, userID); err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	return message, nil
}

func (s *messageService) GetRoomMessages(roomID, userID uint, limit, offset int) ([]*domain.Message, error) {
	// Verify user is participant
	if _, err := s.roomRepo.FindParticipant(roomID, userID); err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	messages, err := s.messageRepo.FindByRoomID(roomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get room messages: %w", err)
	}

	return messages, nil
}

func (s *messageService) Update(messageID, userID uint, req *domain.UpdateMessageRequest) (*domain.Message, error) {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Only sender can edit message
	if message.SenderID != userID {
		return nil, errors.New("only sender can edit message")
	}

	// Can't edit deleted messages
	if message.IsDeleted {
		return nil, errors.New("cannot edit deleted message")
	}

	// Update content
	message.Content = req.Content
	message.IsEdited = true
	now := time.Now()
	message.EditedAt = &now

	if err := s.messageRepo.Update(message); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	// Reload message
	message, err = s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload message: %w", err)
	}

	// Broadcast message edited event
	s.broadcastMessageEvent(message.RoomID, userID, websocket.MessageTypeMessageEdited, message)

	return message, nil
}

func (s *messageService) Delete(messageID, userID uint) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Check if user is sender or room admin
	if message.SenderID != userID {
		participant, err := s.roomRepo.FindParticipant(message.RoomID, userID)
		if err != nil {
			return errors.New("access denied")
		}

		if participant.Role != "admin" {
			room, err := s.roomRepo.FindByID(message.RoomID)
			if err != nil {
				return fmt.Errorf("failed to get room: %w", err)
			}
			if room.CreatorID != userID {
				return errors.New("only sender, admin, or creator can delete message")
			}
		}
	}

	if err := s.messageRepo.SoftDelete(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	// Broadcast message deleted event
	s.broadcastMessageEvent(message.RoomID, userID, websocket.MessageTypeMessageDeleted, map[string]interface{}{
		"message_id": messageID,
		"room_id":    message.RoomID,
	})

	return nil
}

func (s *messageService) AddReaction(messageID, userID uint, req *domain.AddReactionRequest) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Verify user is participant
	if _, err := s.roomRepo.FindParticipant(message.RoomID, userID); err != nil {
		return errors.New("access denied: user is not a participant")
	}

	// Add reaction
	reaction := &domain.MessageReaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     req.Emoji,
	}

	if err := s.messageRepo.AddReaction(reaction); err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	// Broadcast reaction added event
	s.broadcastMessageEvent(message.RoomID, userID, websocket.MessageTypeReactionAdded, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
		"emoji":      req.Emoji,
	})

	return nil
}

func (s *messageService) RemoveReaction(messageID, userID uint, emoji string) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Verify user is participant
	if _, err := s.roomRepo.FindParticipant(message.RoomID, userID); err != nil {
		return errors.New("access denied: user is not a participant")
	}

	if err := s.messageRepo.RemoveReaction(messageID, userID, emoji); err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	// Broadcast reaction removed event
	s.broadcastMessageEvent(message.RoomID, userID, websocket.MessageTypeReactionRemoved, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
		"emoji":      emoji,
	})

	return nil
}

func (s *messageService) MarkAsRead(messageID, userID uint) error {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Verify user is participant
	if _, err := s.roomRepo.FindParticipant(message.RoomID, userID); err != nil {
		return errors.New("access denied: user is not a participant")
	}

	if err := s.messageRepo.MarkAsRead(messageID, userID); err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	// Update room's last_read_at for this user
	if err := s.roomRepo.UpdateLastRead(message.RoomID, userID); err != nil {
		return fmt.Errorf("failed to update last read: %w", err)
	}

	// Broadcast message read event
	s.broadcastMessageEvent(message.RoomID, userID, websocket.MessageTypeMessageRead, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
	})

	return nil
}

func (s *messageService) SendTypingIndicator(roomID, userID uint, isTyping bool) error {
	// Verify user is participant
	if _, err := s.roomRepo.FindParticipant(roomID, userID); err != nil {
		return errors.New("access denied: user is not a participant")
	}

	// Get user info
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Broadcast typing indicator
	indicator := &domain.TypingIndicator{
		RoomID:    roomID,
		UserID:    userID,
		Username:  user.Username,
		IsTyping:  isTyping,
		Timestamp: time.Now(),
	}

	s.broadcastMessageEvent(roomID, userID, websocket.MessageTypeTyping, indicator)

	return nil
}

// Helper methods

func (s *messageService) broadcastMessageEvent(roomID, userID uint, eventType websocket.MessageType, data interface{}) {
	if s.hub != nil {
		message := websocket.NewMessage(eventType, roomID, userID, data)
		s.hub.Broadcast(message)
	}
}
