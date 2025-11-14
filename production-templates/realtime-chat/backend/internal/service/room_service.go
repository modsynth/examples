package service

import (
	"errors"
	"fmt"
	"time"

	"realtime-chat/internal/domain"
	"realtime-chat/internal/repository"
	"realtime-chat/internal/websocket"
)

type RoomService interface {
	Create(creatorID uint, req *domain.CreateRoomRequest) (*domain.Room, error)
	GetByID(roomID, userID uint) (*domain.Room, error)
	GetUserRooms(userID uint) ([]*domain.Room, error)
	Update(roomID, userID uint, req *domain.UpdateRoomRequest) (*domain.Room, error)
	Delete(roomID, userID uint) error
	Archive(roomID, userID uint) error

	// Participant management
	AddParticipant(roomID, requestUserID uint, req *domain.AddParticipantRequest) error
	RemoveParticipant(roomID, participantUserID, requestUserID uint) error
	LeaveRoom(roomID, userID uint) error
	GetParticipants(roomID, userID uint) ([]*domain.Participant, error)

	// Direct message
	GetOrCreateDirectRoom(user1ID, user2ID uint) (*domain.Room, error)

	// Unread count
	GetUnreadCount(roomID, userID uint) (int64, error)
	MarkAsRead(roomID, userID uint) error
}

type roomService struct {
	roomRepo    repository.RoomRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	hub         *websocket.Hub
}

func NewRoomService(
	roomRepo repository.RoomRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	hub *websocket.Hub,
) RoomService {
	return &roomService{
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		hub:         hub,
	}
}

func (s *roomService) Create(creatorID uint, req *domain.CreateRoomRequest) (*domain.Room, error) {
	if req.Name == "" && req.Type != domain.RoomTypeDirect {
		return nil, errors.New("room name is required for non-direct rooms")
	}

	// Verify creator exists
	creator, err := s.userRepo.FindByID(creatorID)
	if err != nil {
		return nil, fmt.Errorf("creator not found: %w", err)
	}

	// For direct rooms, verify exactly 2 participants
	if req.Type == domain.RoomTypeDirect {
		if len(req.UserIDs) != 1 {
			return nil, errors.New("direct room must have exactly one other participant")
		}

		// Check if direct room already exists
		existingRoom, err := s.roomRepo.FindDirectRoom(creatorID, req.UserIDs[0])
		if err != nil {
			return nil, fmt.Errorf("failed to check for existing direct room: %w", err)
		}
		if existingRoom != nil {
			return existingRoom, nil
		}
	}

	// Create room
	room := &domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatorID:   creatorID,
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Add creator as admin participant
	creatorParticipant := &domain.Participant{
		RoomID:   room.ID,
		UserID:   creatorID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}
	if err := s.roomRepo.AddParticipant(creatorParticipant); err != nil {
		return nil, fmt.Errorf("failed to add creator as participant: %w", err)
	}

	// Add other participants
	for _, userID := range req.UserIDs {
		participant := &domain.Participant{
			RoomID:   room.ID,
			UserID:   userID,
			Role:     "member",
			JoinedAt: time.Now(),
		}
		if err := s.roomRepo.AddParticipant(participant); err != nil {
			return nil, fmt.Errorf("failed to add participant: %w", err)
		}

		// Broadcast user joined event
		s.broadcastRoomEvent(room.ID, creatorID, websocket.MessageTypeUserJoined, map[string]interface{}{
			"room_id": room.ID,
			"user_id": userID,
		})
	}

	// Reload room with participants
	return s.roomRepo.FindByID(room.ID)
}

func (s *roomService) GetByID(roomID, userID uint) (*domain.Room, error) {
	// Check if user is participant
	if _, err := s.roomRepo.FindParticipant(roomID, userID); err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	// Load last message
	lastMessage, _ := s.messageRepo.GetLastMessage(roomID)
	room.LastMessage = lastMessage

	// Calculate unread count for each participant
	for i := range room.Participants {
		unreadCount, _ := s.roomRepo.GetUnreadCount(roomID, room.Participants[i].UserID)
		room.Participants[i].UnreadCount = int(unreadCount)
	}

	return room, nil
}

func (s *roomService) GetUserRooms(userID uint) ([]*domain.Room, error) {
	rooms, err := s.roomRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rooms: %w", err)
	}

	// Load last message and unread count for each room
	for i := range rooms {
		lastMessage, _ := s.messageRepo.GetLastMessage(rooms[i].ID)
		rooms[i].LastMessage = lastMessage

		unreadCount, _ := s.roomRepo.GetUnreadCount(rooms[i].ID, userID)
		for j := range rooms[i].Participants {
			if rooms[i].Participants[j].UserID == userID {
				rooms[i].Participants[j].UnreadCount = int(unreadCount)
			}
		}
	}

	return rooms, nil
}

func (s *roomService) Update(roomID, userID uint, req *domain.UpdateRoomRequest) (*domain.Room, error) {
	// Check if user is admin or creator
	participant, err := s.roomRepo.FindParticipant(roomID, userID)
	if err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	if participant.Role != "admin" {
		room, err := s.roomRepo.FindByID(roomID)
		if err != nil {
			return nil, fmt.Errorf("failed to get room: %w", err)
		}
		if room.CreatorID != userID {
			return nil, errors.New("only admin or creator can update room")
		}
	}

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.AvatarURL != "" {
		room.AvatarURL = req.AvatarURL
	}

	if err := s.roomRepo.Update(room); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	// Broadcast room updated event
	s.broadcastRoomEvent(roomID, userID, websocket.MessageTypeRoomUpdated, room)

	return room, nil
}

func (s *roomService) Delete(roomID, userID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Only creator can delete room
	if room.CreatorID != userID {
		return errors.New("only creator can delete room")
	}

	if err := s.roomRepo.Delete(roomID); err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	return nil
}

func (s *roomService) Archive(roomID, userID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Only creator or admin can archive room
	participant, err := s.roomRepo.FindParticipant(roomID, userID)
	if err != nil {
		return errors.New("access denied: user is not a participant")
	}

	if room.CreatorID != userID && participant.Role != "admin" {
		return errors.New("only creator or admin can archive room")
	}

	room.IsArchived = true
	if err := s.roomRepo.Update(room); err != nil {
		return fmt.Errorf("failed to archive room: %w", err)
	}

	return nil
}

func (s *roomService) AddParticipant(roomID, requestUserID uint, req *domain.AddParticipantRequest) error {
	// Check if requester is admin or creator
	participant, err := s.roomRepo.FindParticipant(roomID, requestUserID)
	if err != nil {
		return errors.New("access denied: user is not a participant")
	}

	if participant.Role != "admin" {
		room, err := s.roomRepo.FindByID(roomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}
		if room.CreatorID != requestUserID {
			return errors.New("only admin or creator can add participants")
		}
	}

	// Verify user to add exists
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return fmt.Errorf("user to add not found: %w", err)
	}

	// Check if user is already a participant
	existingParticipant, _ := s.roomRepo.FindParticipant(roomID, req.UserID)
	if existingParticipant != nil {
		return errors.New("user is already a participant")
	}

	// Add participant
	role := req.Role
	if role == "" {
		role = "member"
	}

	newParticipant := &domain.Participant{
		RoomID:   roomID,
		UserID:   req.UserID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := s.roomRepo.AddParticipant(newParticipant); err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	// Broadcast user joined event
	s.broadcastRoomEvent(roomID, requestUserID, websocket.MessageTypeUserJoined, map[string]interface{}{
		"room_id": roomID,
		"user_id": req.UserID,
	})

	return nil
}

func (s *roomService) RemoveParticipant(roomID, participantUserID, requestUserID uint) error {
	// Check if requester is admin or creator
	participant, err := s.roomRepo.FindParticipant(roomID, requestUserID)
	if err != nil {
		return errors.New("access denied: user is not a participant")
	}

	if participant.Role != "admin" {
		room, err := s.roomRepo.FindByID(roomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}
		if room.CreatorID != requestUserID && requestUserID != participantUserID {
			return errors.New("only admin, creator, or the participant themselves can remove")
		}
	}

	if err := s.roomRepo.RemoveParticipant(roomID, participantUserID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	// Broadcast user left event
	s.broadcastRoomEvent(roomID, requestUserID, websocket.MessageTypeUserLeft, map[string]interface{}{
		"room_id": roomID,
		"user_id": participantUserID,
	})

	return nil
}

func (s *roomService) LeaveRoom(roomID, userID uint) error {
	return s.RemoveParticipant(roomID, userID, userID)
}

func (s *roomService) GetParticipants(roomID, userID uint) ([]*domain.Participant, error) {
	// Check if user is participant
	if _, err := s.roomRepo.FindParticipant(roomID, userID); err != nil {
		return nil, errors.New("access denied: user is not a participant")
	}

	participants, err := s.roomRepo.GetParticipants(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return participants, nil
}

func (s *roomService) GetOrCreateDirectRoom(user1ID, user2ID uint) (*domain.Room, error) {
	// Check if direct room already exists
	room, err := s.roomRepo.FindDirectRoom(user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing direct room: %w", err)
	}

	if room != nil {
		return room, nil
	}

	// Create new direct room
	req := &domain.CreateRoomRequest{
		Type:    domain.RoomTypeDirect,
		UserIDs: []uint{user2ID},
	}

	return s.Create(user1ID, req)
}

func (s *roomService) GetUnreadCount(roomID, userID uint) (int64, error) {
	count, err := s.roomRepo.GetUnreadCount(roomID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (s *roomService) MarkAsRead(roomID, userID uint) error {
	if err := s.roomRepo.UpdateLastRead(roomID, userID); err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}
	return nil
}

// Helper methods

func (s *roomService) broadcastRoomEvent(roomID, userID uint, eventType websocket.MessageType, data interface{}) {
	if s.hub != nil {
		message := websocket.NewMessage(eventType, roomID, userID, data)
		s.hub.Broadcast(message)
	}
}
