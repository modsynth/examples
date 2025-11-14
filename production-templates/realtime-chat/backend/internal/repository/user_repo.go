package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"realtime-chat/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	UpdateStatus(userID uint, status domain.UserStatus) error
	UpdateLastSeen(userID uint) error
	Search(query string, limit int) ([]*domain.User, error)
	List(limit, offset int) ([]*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email %s", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with username %s", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateStatus(userID uint, status domain.UserStatus) error {
	if err := r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateLastSeen(userID uint) error {
	if err := r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("last_seen_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}
	return nil
}

func (r *userRepository) Search(query string, limit int) ([]*domain.User, error) {
	var users []*domain.User
	searchPattern := "%" + query + "%"

	err := r.db.Where("username ILIKE ? OR display_name ILIKE ? OR email ILIKE ?",
		searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return users, nil
}

func (r *userRepository) List(limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
