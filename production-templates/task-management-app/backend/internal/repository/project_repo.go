package repository

import (
	"fmt"

	"task-management-app/internal/domain"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *domain.Project) error
	FindByID(id uint) (*domain.Project, error)
	FindByUserID(userID uint) ([]*domain.Project, error)
	Update(project *domain.Project) error
	Delete(id uint) error
	AddMember(member *domain.ProjectMember) error
	RemoveMember(projectID, userID uint) error
	UpdateMember(member *domain.ProjectMember) error
	GetMember(projectID, userID uint) (*domain.ProjectMember, error)
	GetMembers(projectID uint) ([]domain.ProjectMember, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *domain.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (r *projectRepository) FindByID(id uint) (*domain.Project, error) {
	var project domain.Project
	err := r.db.Preload("Owner").Preload("Members.User").Preload("Boards").First(&project, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find project: %w", err)
	}
	return &project, nil
}

func (r *projectRepository) FindByUserID(userID uint) ([]*domain.Project, error) {
	var projects []*domain.Project

	err := r.db.
		Joins("LEFT JOIN project_members ON projects.id = project_members.project_id").
		Where("projects.owner_id = ? OR project_members.user_id = ?", userID, userID).
		Group("projects.id").
		Preload("Owner").
		Preload("Members.User").
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find projects for user: %w", err)
	}
	return projects, nil
}

func (r *projectRepository) Update(project *domain.Project) error {
	if err := r.db.Save(project).Error; err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

func (r *projectRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Project{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

func (r *projectRepository) AddMember(member *domain.ProjectMember) error {
	if err := r.db.Create(member).Error; err != nil {
		return fmt.Errorf("failed to add project member: %w", err)
	}
	return nil
}

func (r *projectRepository) RemoveMember(projectID, userID uint) error {
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&domain.ProjectMember{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove project member: %w", err)
	}
	return nil
}

func (r *projectRepository) UpdateMember(member *domain.ProjectMember) error {
	if err := r.db.Save(member).Error; err != nil {
		return fmt.Errorf("failed to update project member: %w", err)
	}
	return nil
}

func (r *projectRepository) GetMember(projectID, userID uint) (*domain.ProjectMember, error) {
	var member domain.ProjectMember
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Preload("User").First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project member not found")
		}
		return nil, fmt.Errorf("failed to get project member: %w", err)
	}
	return &member, nil
}

func (r *projectRepository) GetMembers(projectID uint) ([]domain.ProjectMember, error) {
	var members []domain.ProjectMember
	err := r.db.Where("project_id = ?", projectID).Preload("User").Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get project members: %w", err)
	}
	return members, nil
}
