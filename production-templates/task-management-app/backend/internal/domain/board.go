package domain

import "time"

type Board struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Position  int       `json:"position" gorm:"not null;default:0"`
	Tasks     []Task    `json:"tasks,omitempty" gorm:"foreignKey:BoardID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBoardRequest struct {
	Name     string `json:"name" binding:"required"`
	Position int    `json:"position"`
}

type UpdateBoardRequest struct {
	Name     string `json:"name"`
	Position *int   `json:"position"`
}
