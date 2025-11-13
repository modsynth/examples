package repository

import (
	"fmt"

	"task-management-app/internal/domain"
	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(board *domain.Board) error
	FindByID(id uint) (*domain.Board, error)
	FindByProjectID(projectID uint) ([]*domain.Board, error)
	Update(board *domain.Board) error
	Delete(id uint) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *domain.Board) error {
	if err := r.db.Create(board).Error; err != nil {
		return fmt.Errorf("failed to create board: %w", err)
	}
	return nil
}

func (r *boardRepository) FindByID(id uint) (*domain.Board, error) {
	var board domain.Board
	err := r.db.Preload("Tasks").First(&board, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("board not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find board: %w", err)
	}
	return &board, nil
}

func (r *boardRepository) FindByProjectID(projectID uint) ([]*domain.Board, error) {
	var boards []*domain.Board
	err := r.db.Where("project_id = ?", projectID).
		Order("position ASC").
		Preload("Tasks").
		Find(&boards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find boards by project: %w", err)
	}
	return boards, nil
}

func (r *boardRepository) Update(board *domain.Board) error {
	if err := r.db.Save(board).Error; err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}
	return nil
}

func (r *boardRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Board{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}
	return nil
}
