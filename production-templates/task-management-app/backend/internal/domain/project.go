package domain

import "time"

type ProjectRole string

const (
	ProjectRoleOwner  ProjectRole = "owner"
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleMember ProjectRole = "member"
	ProjectRoleViewer ProjectRole = "viewer"
)

type Project struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description"`
	Icon        string           `json:"icon"`
	Color       string           `json:"color"`
	OwnerID     uint             `json:"owner_id" gorm:"not null"`
	Owner       *User            `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Members     []ProjectMember  `json:"members,omitempty" gorm:"foreignKey:ProjectID"`
	Boards      []Board          `json:"boards,omitempty" gorm:"foreignKey:ProjectID"`
	IsArchived  bool             `json:"is_archived" gorm:"not null;default:false"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ProjectMember struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	ProjectID uint        `json:"project_id" gorm:"not null;uniqueIndex:idx_project_user"`
	UserID    uint        `json:"user_id" gorm:"not null;uniqueIndex:idx_project_user"`
	User      *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role      ProjectRole `json:"role" gorm:"not null;default:'member'"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type AddMemberRequest struct {
	UserID uint        `json:"user_id" binding:"required"`
	Role   ProjectRole `json:"role" binding:"required"`
}

type UpdateMemberRoleRequest struct {
	Role ProjectRole `json:"role" binding:"required"`
}
