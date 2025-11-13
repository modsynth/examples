package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	FindByID(id uint) (*domain.Order, error)
	FindByOrderNumber(orderNumber string) (*domain.Order, error)
	FindByUserID(userID uint, page, limit int) ([]*domain.Order, int64, error)
	Update(order *domain.Order) error
	UpdateStatus(orderID uint, status domain.OrderStatus) error
	UpdatePaymentStatus(orderID uint, status domain.PaymentStatus) error
	List(query *domain.OrderListQuery) ([]*domain.Order, int64, error)
	GenerateOrderNumber() (string, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("User").Preload("Items").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("User").Preload("Items").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint, page, limit int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	db := r.db.Model(&domain.Order{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(orderID uint, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("status", status).
		Error
}

func (r *orderRepository) UpdatePaymentStatus(orderID uint, status domain.PaymentStatus) error {
	return r.db.Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", status).
		Error
}

func (r *orderRepository) List(query *domain.OrderListQuery) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	offset := (query.Page - 1) * query.Limit

	db := r.db.Model(&domain.Order{})

	// Apply filters
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if query.PaymentStatus != nil {
		db = db.Where("payment_status = ?", *query.PaymentStatus)
	}

	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if query.OrderNumber != "" {
		db = db.Where("order_number = ?", query.OrderNumber)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch results
	err := db.Preload("User").
		Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(query.Limit).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) GenerateOrderNumber() (string, error) {
	// Format: ORD-YYYYMMDD-XXXXXX (e.g., ORD-20231114-000001)
	now := time.Now()
	prefix := fmt.Sprintf("ORD-%s", now.Format("20060102"))

	var count int64
	err := r.db.Model(&domain.Order{}).
		Where("order_number LIKE ?", prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	orderNumber := fmt.Sprintf("%s-%06d", prefix, count+1)
	return orderNumber, nil
}
