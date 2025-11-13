package domain

import "time"

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

type Task struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	BoardID     uint            `json:"board_id" gorm:"not null"`
	Board       *Board          `json:"board,omitempty" gorm:"foreignKey:BoardID"`
	Title       string          `json:"title" gorm:"not null"`
	Description string          `json:"description"`
	Position    int             `json:"position" gorm:"not null;default:0"`
	Priority    TaskPriority    `json:"priority" gorm:"not null;default:'medium'"`
	DueDate     *time.Time      `json:"due_date"`
	CreatorID   uint            `json:"creator_id" gorm:"not null"`
	Creator     *User           `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	AssigneeID  *uint           `json:"assignee_id"`
	Assignee    *User           `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	Labels      []Label         `json:"labels,omitempty" gorm:"many2many:task_labels"`
	Comments    []Comment       `json:"comments,omitempty" gorm:"foreignKey:TaskID"`
	Attachments []Attachment    `json:"attachments,omitempty" gorm:"foreignKey:TaskID"`
	Checklist   []ChecklistItem `json:"checklist,omitempty" gorm:"foreignKey:TaskID"`
	IsCompleted bool            `json:"is_completed" gorm:"not null;default:false"`
	CompletedAt *time.Time      `json:"completed_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Label struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Color     string    `json:"color" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Attachment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Filename  string    `json:"filename" gorm:"not null"`
	FileURL   string    `json:"file_url" gorm:"not null"`
	FileSize  int64     `json:"file_size" gorm:"not null"`
	MimeType  string    `json:"mime_type" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

type ChecklistItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TaskID      uint      `json:"task_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	IsCompleted bool      `json:"is_completed" gorm:"not null;default:false"`
	Position    int       `json:"position" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	AssigneeID  *uint        `json:"assignee_id"`
	LabelIDs    []uint       `json:"label_ids"`
}

type UpdateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	AssigneeID  *uint        `json:"assignee_id"`
	IsCompleted *bool        `json:"is_completed"`
}

type MoveTaskRequest struct {
	BoardID  uint `json:"board_id" binding:"required"`
	Position int  `json:"position" binding:"gte=0"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateChecklistItemRequest struct {
	Title    string `json:"title" binding:"required"`
	Position int    `json:"position"`
}

type UpdateChecklistItemRequest struct {
	Title       string `json:"title"`
	IsCompleted *bool  `json:"is_completed"`
}
