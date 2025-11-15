package service

import (
	"errors"
	"fmt"

	"task-management-app/internal/domain"
	"task-management-app/internal/repository"
	"task-management-app/internal/websocket"
)

type BoardService interface {
	Create(projectID, userID uint, req *domain.CreateBoardRequest) (*domain.Board, error)
	GetByID(boardID, userID uint) (*domain.Board, error)
	Update(boardID, userID uint, req *domain.UpdateBoardRequest) (*domain.Board, error)
	Delete(boardID, userID uint) error
	ListByProject(projectID, userID uint) ([]*domain.Board, error)
}

type boardService struct {
	boardRepo   repository.BoardRepository
	projectRepo repository.ProjectRepository
	hub         *websocket.Hub
}

func NewBoardService(
	boardRepo repository.BoardRepository,
	projectRepo repository.ProjectRepository,
	hub *websocket.Hub,
) BoardService {
	return &boardService{
		boardRepo:   boardRepo,
		projectRepo: projectRepo,
		hub:         hub,
	}
}

func (s *boardService) Create(projectID, userID uint, req *domain.CreateBoardRequest) (*domain.Board, error) {
	if req.Name == "" {
		return nil, errors.New("board name is required")
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(projectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	// Get next position for the board
	boards, _ := s.boardRepo.FindByProjectID(projectID)
	position := req.Position
	if position == 0 {
		position = len(boards)
	}

	board := &domain.Board{
		ProjectID: projectID,
		Name:      req.Name,
		Position:  position,
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastBoardEvent(projectID, userID, "BOARD_CREATED", board)

	return board, nil
}

func (s *boardService) GetByID(boardID, userID uint) (*domain.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleViewer); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) Update(boardID, userID uint, req *domain.UpdateBoardRequest) (*domain.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		board.Name = req.Name
	}
	if req.Position != nil {
		board.Position = *req.Position
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastBoardEvent(board.ProjectID, userID, "BOARD_UPDATED", board)

	return board, nil
}

func (s *boardService) Delete(boardID, userID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleAdmin); err != nil {
		return err
	}

	if err := s.boardRepo.Delete(boardID); err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastBoardEvent(board.ProjectID, userID, "BOARD_DELETED", map[string]interface{}{
		"id":         boardID,
		"project_id": board.ProjectID,
	})

	return nil
}

func (s *boardService) ListByProject(projectID, userID uint) ([]*domain.Board, error) {
	// Check if user has access to the project
	if err := s.checkProjectAccess(projectID, userID, domain.ProjectRoleViewer); err != nil {
		return nil, err
	}

	boards, err := s.boardRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list boards: %w", err)
	}

	return boards, nil
}

// Helper methods

func (s *boardService) checkProjectAccess(projectID, userID uint, requiredRole domain.ProjectRole) error {
	member, err := s.projectRepo.GetMember(projectID, userID)
	if err != nil {
		return errors.New("access denied: user is not a member of this project")
	}

	// Check role hierarchy
	roleHierarchy := map[domain.ProjectRole]int{
		domain.ProjectRoleViewer: 1,
		domain.ProjectRoleMember: 2,
		domain.ProjectRoleAdmin:  3,
		domain.ProjectRoleOwner:  4,
	}

	if roleHierarchy[member.Role] < roleHierarchy[requiredRole] {
		return fmt.Errorf("insufficient permissions: required %s role", requiredRole)
	}

	return nil
}

func (s *boardService) broadcastBoardEvent(projectID, userID uint, eventType string, data interface{}) {
	if s.hub != nil {
		message := &websocket.Message{
			Type:      websocket.MessageType(eventType),
			ProjectID: projectID,
			UserID:    userID,
			Payload:   data,
		}
		s.hub.Broadcast(message)
	}
}
