package service

import (
	"errors"
	"fmt"
	"time"

	"task-management-app/internal/domain"
	"task-management-app/internal/repository"
	"task-management-app/internal/websocket"
)

type TaskService interface {
	Create(boardID, userID uint, req *domain.CreateTaskRequest) (*domain.Task, error)
	GetByID(taskID, userID uint) (*domain.Task, error)
	Update(taskID, userID uint, req *domain.UpdateTaskRequest) (*domain.Task, error)
	Delete(taskID, userID uint) error
	Move(taskID, userID uint, req *domain.MoveTaskRequest) error
	ListByBoard(boardID, userID uint) ([]*domain.Task, error)

	AddComment(taskID, userID uint, req *domain.CreateCommentRequest) (*domain.Comment, error)
	DeleteComment(commentID, userID uint) error

	AddChecklistItem(taskID, userID uint, req *domain.CreateChecklistItemRequest) (*domain.ChecklistItem, error)
	UpdateChecklistItem(itemID, userID uint, req *domain.UpdateChecklistItemRequest) (*domain.ChecklistItem, error)
	DeleteChecklistItem(itemID, userID uint) error

	AssignLabels(taskID, userID uint, labelIDs []uint) error
}

type taskService struct {
	taskRepo    repository.TaskRepository
	boardRepo   repository.BoardRepository
	projectRepo repository.ProjectRepository
	hub         *websocket.Hub
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	boardRepo repository.BoardRepository,
	projectRepo repository.ProjectRepository,
	hub *websocket.Hub,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		boardRepo:   boardRepo,
		projectRepo: projectRepo,
		hub:         hub,
	}
}

func (s *taskService) Create(boardID, userID uint, req *domain.CreateTaskRequest) (*domain.Task, error) {
	if req.Title == "" {
		return nil, errors.New("task title is required")
	}

	// Get board to check access and get project ID
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	// Get next position for the task
	tasks, _ := s.taskRepo.FindByBoardID(boardID)
	position := len(tasks)

	// Set default priority
	priority := req.Priority
	if priority == "" {
		priority = domain.PriorityMedium
	}

	task := &domain.Task{
		BoardID:     boardID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		DueDate:     req.DueDate,
		AssigneeID:  req.AssigneeID,
		CreatorID:   userID,
		Position:    position,
		IsCompleted: false,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Assign labels if provided
	if len(req.LabelIDs) > 0 {
		if err := s.taskRepo.AssignLabels(task.ID, req.LabelIDs); err != nil {
			return nil, fmt.Errorf("failed to assign labels: %w", err)
		}
	}

	// Reload task with all relations
	task, err = s.taskRepo.FindByID(task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload task: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "TASK_CREATED", task)

	return task, nil
}

func (s *taskService) GetByID(taskID, userID uint) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Get board to check access
	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleViewer); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) Update(taskID, userID uint, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Get board to check access
	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.IsCompleted != nil {
		task.IsCompleted = *req.IsCompleted
		if *req.IsCompleted {
			now := time.Now()
			task.CompletedAt = &now
		} else {
			task.CompletedAt = nil
		}
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Reload task with all relations
	task, err = s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload task: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "TASK_UPDATED", task)

	return task, nil
}

func (s *taskService) Delete(taskID, userID uint) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Get board to check access
	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return err
	}

	if err := s.taskRepo.Delete(taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "TASK_DELETED", map[string]interface{}{
		"id":       taskID,
		"board_id": task.BoardID,
	})

	return nil
}

func (s *taskService) Move(taskID, userID uint, req *domain.MoveTaskRequest) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Get source board
	sourceBoard, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return fmt.Errorf("source board not found: %w", err)
	}

	// Get target board
	targetBoard, err := s.boardRepo.FindByID(req.BoardID)
	if err != nil {
		return fmt.Errorf("target board not found: %w", err)
	}

	// Both boards must be in the same project
	if sourceBoard.ProjectID != targetBoard.ProjectID {
		return errors.New("cannot move task between different projects")
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(sourceBoard.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return err
	}

	if err := s.taskRepo.Move(taskID, req.BoardID, req.Position); err != nil {
		return fmt.Errorf("failed to move task: %w", err)
	}

	// Reload task with all relations
	task, err = s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to reload task: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(sourceBoard.ProjectID, userID, "TASK_MOVED", task)

	return nil
}

func (s *taskService) ListByBoard(boardID, userID uint) ([]*domain.Task, error) {
	// Get board to check access
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleViewer); err != nil {
		return nil, err
	}

	tasks, err := s.taskRepo.FindByBoardID(boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) AddComment(taskID, userID uint, req *domain.CreateCommentRequest) (*domain.Comment, error) {
	if req.Content == "" {
		return nil, errors.New("comment content is required")
	}

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Get board to check access
	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	// Check if user has access to the project
	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	comment := &domain.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.taskRepo.AddComment(comment); err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "COMMENT_ADDED", comment)

	return comment, nil
}

func (s *taskService) DeleteComment(commentID, userID uint) error {
	comment, err := s.taskRepo.GetComment(commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Only the comment author can delete it
	if comment.UserID != userID {
		return errors.New("only comment author can delete the comment")
	}

	task, err := s.taskRepo.FindByID(comment.TaskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}

	if err := s.taskRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "COMMENT_DELETED", map[string]interface{}{
		"id":      commentID,
		"task_id": task.ID,
	})

	return nil
}

func (s *taskService) AddChecklistItem(taskID, userID uint, req *domain.CreateChecklistItemRequest) (*domain.ChecklistItem, error) {
	if req.Title == "" {
		return nil, errors.New("checklist item title is required")
	}

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	item := &domain.ChecklistItem{
		TaskID:   taskID,
		Title:    req.Title,
		Position: req.Position,
	}

	if err := s.taskRepo.AddChecklistItem(item); err != nil {
		return nil, fmt.Errorf("failed to add checklist item: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "CHECKLIST_ITEM_ADDED", item)

	return item, nil
}

func (s *taskService) UpdateChecklistItem(itemID, userID uint, req *domain.UpdateChecklistItemRequest) (*domain.ChecklistItem, error) {
	item, err := s.taskRepo.GetChecklistItem(itemID)
	if err != nil {
		return nil, fmt.Errorf("checklist item not found: %w", err)
	}

	task, err := s.taskRepo.FindByID(item.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return nil, fmt.Errorf("board not found: %w", err)
	}

	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return nil, err
	}

	if req.Title != "" {
		item.Title = req.Title
	}
	if req.IsCompleted != nil {
		item.IsCompleted = *req.IsCompleted
	}

	if err := s.taskRepo.UpdateChecklistItem(item); err != nil {
		return nil, fmt.Errorf("failed to update checklist item: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "CHECKLIST_ITEM_UPDATED", item)

	return item, nil
}

func (s *taskService) DeleteChecklistItem(itemID, userID uint) error {
	item, err := s.taskRepo.GetChecklistItem(itemID)
	if err != nil {
		return fmt.Errorf("checklist item not found: %w", err)
	}

	task, err := s.taskRepo.FindByID(item.TaskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}

	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return err
	}

	if err := s.taskRepo.DeleteChecklistItem(itemID); err != nil {
		return fmt.Errorf("failed to delete checklist item: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "CHECKLIST_ITEM_DELETED", map[string]interface{}{
		"id":      itemID,
		"task_id": task.ID,
	})

	return nil
}

func (s *taskService) AssignLabels(taskID, userID uint, labelIDs []uint) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	board, err := s.boardRepo.FindByID(task.BoardID)
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}

	if err := s.checkProjectAccess(board.ProjectID, userID, domain.ProjectRoleMember); err != nil {
		return err
	}

	if err := s.taskRepo.AssignLabels(taskID, labelIDs); err != nil {
		return fmt.Errorf("failed to assign labels: %w", err)
	}

	// Reload task to get updated labels
	task, err = s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to reload task: %w", err)
	}

	// Broadcast via WebSocket
	s.broadcastTaskEvent(board.ProjectID, userID, "TASK_LABELS_UPDATED", task)

	return nil
}

// Helper methods

func (s *taskService) checkProjectAccess(projectID, userID uint, requiredRole domain.ProjectRole) error {
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

func (s *taskService) broadcastTaskEvent(projectID, userID uint, eventType string, data interface{}) {
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
