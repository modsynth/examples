package repository

import (
	"fmt"

	"task-management-app/internal/domain"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id uint) (*domain.Task, error)
	FindByBoardID(boardID uint) ([]*domain.Task, error)
	FindByProjectID(projectID uint) ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error
	Move(taskID, boardID uint, position int) error
	AddComment(comment *domain.Comment) error
	GetComment(commentID uint) (*domain.Comment, error)
	DeleteComment(commentID uint) error
	GetComments(taskID uint) ([]*domain.Comment, error)
	AddAttachment(attachment *domain.Attachment) error
	GetAttachments(taskID uint) ([]*domain.Attachment, error)
	AddChecklistItem(item *domain.ChecklistItem) error
	GetChecklistItem(itemID uint) (*domain.ChecklistItem, error)
	UpdateChecklistItem(item *domain.ChecklistItem) error
	DeleteChecklistItem(id uint) error
	AssignLabels(taskID uint, labelIDs []uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (r *taskRepository) FindByID(id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.
		Preload("Board").
		Preload("Creator").
		Preload("Assignee").
		Preload("Labels").
		Preload("Comments.User").
		Preload("Attachments.User").
		Preload("Checklist").
		First(&task, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	return &task, nil
}

func (r *taskRepository) FindByBoardID(boardID uint) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.Where("board_id = ?", boardID).
		Preload("Creator").
		Preload("Assignee").
		Preload("Labels").
		Order("position ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find tasks by board: %w", err)
	}
	return tasks, nil
}

func (r *taskRepository) FindByProjectID(projectID uint) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.
		Joins("JOIN boards ON tasks.board_id = boards.id").
		Where("boards.project_id = ?", projectID).
		Preload("Board").
		Preload("Creator").
		Preload("Assignee").
		Preload("Labels").
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find tasks by project: %w", err)
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	if err := r.db.Save(task).Error; err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func (r *taskRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Task{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (r *taskRepository) Move(taskID, boardID uint, position int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update task board and position
		if err := tx.Model(&domain.Task{}).
			Where("id = ?", taskID).
			Updates(map[string]interface{}{
				"board_id": boardID,
				"position": position,
			}).Error; err != nil {
			return fmt.Errorf("failed to move task: %w", err)
		}

		// Reorder other tasks in the target board
		if err := tx.Exec(
			"UPDATE tasks SET position = position + 1 WHERE board_id = ? AND id != ? AND position >= ?",
			boardID, taskID, position,
		).Error; err != nil {
			return fmt.Errorf("failed to reorder tasks: %w", err)
		}

		return nil
	})
}

func (r *taskRepository) AddComment(comment *domain.Comment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	// Reload comment with user
	return r.db.Preload("User").First(comment, comment.ID).Error
}

func (r *taskRepository) GetComment(commentID uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.Preload("User").First(&comment, commentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment not found with id %d", commentID)
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

func (r *taskRepository) DeleteComment(commentID uint) error {
	if err := r.db.Delete(&domain.Comment{}, commentID).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (r *taskRepository) GetComments(taskID uint) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	err := r.db.Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, nil
}

func (r *taskRepository) AddAttachment(attachment *domain.Attachment) error {
	if err := r.db.Create(attachment).Error; err != nil {
		return fmt.Errorf("failed to add attachment: %w", err)
	}
	return nil
}

func (r *taskRepository) GetAttachments(taskID uint) ([]*domain.Attachment, error) {
	var attachments []*domain.Attachment
	err := r.db.Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Find(&attachments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}
	return attachments, nil
}

func (r *taskRepository) AddChecklistItem(item *domain.ChecklistItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return fmt.Errorf("failed to add checklist item: %w", err)
	}
	return nil
}

func (r *taskRepository) GetChecklistItem(itemID uint) (*domain.ChecklistItem, error) {
	var item domain.ChecklistItem
	err := r.db.First(&item, itemID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("checklist item not found with id %d", itemID)
		}
		return nil, fmt.Errorf("failed to get checklist item: %w", err)
	}
	return &item, nil
}

func (r *taskRepository) UpdateChecklistItem(item *domain.ChecklistItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return fmt.Errorf("failed to update checklist item: %w", err)
	}
	return nil
}

func (r *taskRepository) DeleteChecklistItem(id uint) error {
	if err := r.db.Delete(&domain.ChecklistItem{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete checklist item: %w", err)
	}
	return nil
}

func (r *taskRepository) AssignLabels(taskID uint, labelIDs []uint) error {
	var task domain.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Clear existing labels
	if err := r.db.Model(&task).Association("Labels").Clear(); err != nil {
		return fmt.Errorf("failed to clear labels: %w", err)
	}

	// Assign new labels
	if len(labelIDs) > 0 {
		var labels []domain.Label
		if err := r.db.Find(&labels, labelIDs).Error; err != nil {
			return fmt.Errorf("failed to find labels: %w", err)
		}
		if err := r.db.Model(&task).Association("Labels").Append(&labels); err != nil {
			return fmt.Errorf("failed to assign labels: %w", err)
		}
	}

	return nil
}
