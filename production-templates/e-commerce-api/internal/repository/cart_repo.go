package repository

import (
	"errors"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserID(userID uint) (*domain.Cart, error)
	CreateCart(cart *domain.Cart) error
	AddItem(item *domain.CartItem) error
	UpdateItem(item *domain.CartItem) error
	RemoveItem(itemID uint) error
	ClearCart(userID uint) error
	GetCartWithItems(userID uint) (*domain.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) CreateCart(cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) AddItem(item *domain.CartItem) error {
	// Check if item already exists
	var existingItem domain.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).First(&existingItem).Error

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += item.Quantity
		return r.db.Save(&existingItem).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Item doesn't exist, create new
		return r.db.Create(item).Error
	}

	return err
}

func (r *cartRepository) UpdateItem(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) RemoveItem(itemID uint) error {
	return r.db.Delete(&domain.CartItem{}, itemID).Error
}

func (r *cartRepository) ClearCart(userID uint) error {
	cart, err := r.FindByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error
}

func (r *cartRepository) GetCartWithItems(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items.Product.Images").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new cart if doesn't exist
			cart = domain.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}
	return &cart, nil
}
