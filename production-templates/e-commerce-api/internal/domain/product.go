package domain

import "time"

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	CategoryID     *uint           `json:"category_id"`
	Category       *Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Name           string          `json:"name" gorm:"not null"`
	Slug           string          `json:"slug" gorm:"uniqueIndex;not null"`
	Description    string          `json:"description"`
	Price          float64         `json:"price" gorm:"not null"`
	ComparePrice   *float64        `json:"compare_price,omitempty"`
	CostPrice      *float64        `json:"cost_price,omitempty"`
	SKU            string          `json:"sku" gorm:"uniqueIndex"`
	Barcode        string          `json:"barcode"`
	StockQuantity  int             `json:"stock_quantity" gorm:"not null;default:0"`
	TrackInventory bool            `json:"track_inventory" gorm:"not null;default:true"`
	Weight         *float64        `json:"weight,omitempty"`
	IsActive       bool            `json:"is_active" gorm:"not null;default:true"`
	Featured       bool            `json:"featured" gorm:"not null;default:false"`
	Images         []ProductImage  `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ProductImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	URL       string    `json:"url" gorm:"not null"`
	AltText   string    `json:"alt_text"`
	Position  int       `json:"position" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	CategoryID     *uint    `json:"category_id"`
	Name           string   `json:"name" binding:"required"`
	Slug           string   `json:"slug" binding:"required"`
	Description    string   `json:"description"`
	Price          float64  `json:"price" binding:"required,gt=0"`
	ComparePrice   *float64 `json:"compare_price"`
	CostPrice      *float64 `json:"cost_price"`
	SKU            string   `json:"sku"`
	Barcode        string   `json:"barcode"`
	StockQuantity  int      `json:"stock_quantity" binding:"gte=0"`
	TrackInventory bool     `json:"track_inventory"`
	Weight         *float64 `json:"weight"`
	IsActive       bool     `json:"is_active"`
	Featured       bool     `json:"featured"`
}

type UpdateProductRequest struct {
	CategoryID     *uint    `json:"category_id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Price          *float64 `json:"price" binding:"omitempty,gt=0"`
	ComparePrice   *float64 `json:"compare_price"`
	CostPrice      *float64 `json:"cost_price"`
	SKU            string   `json:"sku"`
	Barcode        string   `json:"barcode"`
	StockQuantity  *int     `json:"stock_quantity" binding:"omitempty,gte=0"`
	TrackInventory *bool    `json:"track_inventory"`
	Weight         *float64 `json:"weight"`
	IsActive       *bool    `json:"is_active"`
	Featured       *bool    `json:"featured"`
}

type ProductListQuery struct {
	Page       int     `form:"page" binding:"omitempty,gte=1"`
	Limit      int     `form:"limit" binding:"omitempty,gte=1,lte=100"`
	CategoryID *uint   `form:"category_id"`
	Search     string  `form:"search"`
	MinPrice   *float64 `form:"min_price" binding:"omitempty,gte=0"`
	MaxPrice   *float64 `form:"max_price" binding:"omitempty,gte=0"`
	IsActive   *bool   `form:"is_active"`
	Featured   *bool   `form:"featured"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
}
