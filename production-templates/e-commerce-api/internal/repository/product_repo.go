package repository

import (
	"errors"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	FindByID(id uint) (*domain.Product, error)
	FindBySlug(slug string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	List(query *domain.ProductListQuery) ([]*domain.Product, int64, error)
	DecrementStock(productID uint, quantity int) error
	IncrementStock(productID uint, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Preload("Images").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Preload("Images").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) List(query *domain.ProductListQuery) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	offset := (query.Page - 1) * query.Limit

	// Build query
	db := r.db.Model(&domain.Product{})

	// Apply filters
	if query.CategoryID != nil {
		db = db.Where("category_id = ?", *query.CategoryID)
	}

	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}

	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.Featured != nil {
		db = db.Where("featured = ?", *query.Featured)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	sortOrder := "DESC"
	if query.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Fetch results
	err := db.Preload("Category").
		Preload("Images").
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(query.Limit).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) DecrementStock(productID uint, quantity int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, quantity).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", quantity)).
		Error
}

func (r *productRepository) IncrementStock(productID uint, quantity int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).
		Error
}
