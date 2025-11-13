package service

import (
	"errors"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/repository"
)

type CartService interface {
	GetCart(userID uint) (*domain.CartWithSummary, error)
	AddToCart(userID uint, req *domain.AddToCartRequest) error
	UpdateCartItem(userID, itemID uint, req *domain.UpdateCartItemRequest) error
	RemoveFromCart(userID, itemID uint) error
	ClearCart(userID uint) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) GetCart(userID uint) (*domain.CartWithSummary, error) {
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return nil, err
	}

	// Calculate summary
	subtotal := 0.0
	itemsCount := 0

	for _, item := range cart.Items {
		subtotal += item.Price * float64(item.Quantity)
		itemsCount += item.Quantity
	}

	return &domain.CartWithSummary{
		Cart:       cart,
		Subtotal:   subtotal,
		ItemsCount: itemsCount,
	}, nil
}

func (s *cartService) AddToCart(userID uint, req *domain.AddToCartRequest) error {
	// Get or create cart
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return err
	}

	// Check if product exists and has enough stock
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if !product.IsActive {
		return errors.New("product is not available")
	}

	if product.TrackInventory && product.StockQuantity < req.Quantity {
		return errors.New("insufficient stock")
	}

	// Add item to cart
	cartItem := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price,
	}

	return s.cartRepo.AddItem(cartItem)
}

func (s *cartService) UpdateCartItem(userID, itemID uint, req *domain.UpdateCartItemRequest) error {
	// Get cart to verify ownership
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return err
	}

	// Find cart item
	var cartItem *domain.CartItem
	for i, item := range cart.Items {
		if item.ID == itemID {
			cartItem = &cart.Items[i]
			break
		}
	}

	if cartItem == nil {
		return errors.New("cart item not found")
	}

	// Check stock
	product, err := s.productRepo.FindByID(cartItem.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.TrackInventory && product.StockQuantity < req.Quantity {
		return errors.New("insufficient stock")
	}

	// Update quantity
	cartItem.Quantity = req.Quantity

	return s.cartRepo.UpdateItem(cartItem)
}

func (s *cartService) RemoveFromCart(userID, itemID uint) error {
	// Get cart to verify ownership
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return err
	}

	// Verify item belongs to user's cart
	found := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("cart item not found")
	}

	return s.cartRepo.RemoveItem(itemID)
}

func (s *cartService) ClearCart(userID uint) error {
	return s.cartRepo.ClearCart(userID)
}
