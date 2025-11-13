package service

import (
	"errors"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/repository"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(userID uint, req *domain.CreateOrderRequest) (*domain.Order, error)
	GetOrderByID(userID, orderID uint) (*domain.Order, error)
	GetOrderByOrderNumber(userID uint, orderNumber string) (*domain.Order, error)
	GetUserOrders(userID uint, page, limit int) ([]*domain.Order, int64, error)
	CancelOrder(userID, orderID uint) error
	// Admin methods
	GetAllOrders(query *domain.OrderListQuery) ([]*domain.Order, int64, error)
	UpdateOrderStatus(orderID uint, status domain.OrderStatus) error
}

type orderService struct {
	db          *gorm.DB
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
) OrderService {
	return &orderService{
		db:          db,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(userID uint, req *domain.CreateOrderRequest) (*domain.Order, error) {
	var order *domain.Order

	// Use transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get cart with items
		cart, err := s.cartRepo.GetCartWithItems(userID)
		if err != nil {
			return err
		}

		if len(cart.Items) == 0 {
			return errors.New("cart is empty")
		}

		// Calculate totals and create order items
		var orderItems []domain.OrderItem
		subtotal := 0.0

		for _, cartItem := range cart.Items {
			// Check stock availability
			product, err := s.productRepo.FindByID(cartItem.ProductID)
			if err != nil {
				return errors.New("product not found: " + err.Error())
			}

			if product.TrackInventory && product.StockQuantity < cartItem.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}

			// Create order item
			itemSubtotal := cartItem.Price * float64(cartItem.Quantity)
			orderItem := domain.OrderItem{
				ProductID:   cartItem.ProductID,
				ProductName: product.Name,
				ProductSKU:  product.SKU,
				Quantity:    cartItem.Quantity,
				Price:       cartItem.Price,
				Subtotal:    itemSubtotal,
			}

			orderItems = append(orderItems, orderItem)
			subtotal += itemSubtotal

			// Decrement stock
			if product.TrackInventory {
				if err := s.productRepo.DecrementStock(cartItem.ProductID, cartItem.Quantity); err != nil {
					return errors.New("failed to decrement stock")
				}
			}
		}

		// Calculate tax and shipping (simplified)
		tax := subtotal * 0.1 // 10% tax
		shipping := 10.0      // Flat shipping rate
		total := subtotal + tax + shipping

		// Generate order number
		orderNumber, err := s.orderRepo.GenerateOrderNumber()
		if err != nil {
			return errors.New("failed to generate order number")
		}

		// Create order
		order = &domain.Order{
			UserID:               userID,
			OrderNumber:          orderNumber,
			Status:               domain.OrderStatusPending,
			Subtotal:             subtotal,
			Tax:                  tax,
			Shipping:             shipping,
			Total:                total,
			Currency:             "USD",
			PaymentStatus:        domain.PaymentStatusPending,
			PaymentMethod:        req.PaymentMethod,
			ShippingAddressLine1: req.ShippingAddress.Line1,
			ShippingAddressLine2: req.ShippingAddress.Line2,
			ShippingCity:         req.ShippingAddress.City,
			ShippingState:        req.ShippingAddress.State,
			ShippingPostalCode:   req.ShippingAddress.PostalCode,
			ShippingCountry:      req.ShippingAddress.Country,
			Notes:                req.Notes,
			Items:                orderItems,
		}

		if err := s.orderRepo.Create(order); err != nil {
			return errors.New("failed to create order")
		}

		// Clear cart
		if err := s.cartRepo.ClearCart(userID); err != nil {
			return errors.New("failed to clear cart")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrderByID(userID, orderID uint) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (s *orderService) GetOrderByOrderNumber(userID uint, orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.FindByOrderNumber(orderNumber)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID uint, page, limit int) ([]*domain.Order, int64, error) {
	return s.orderRepo.FindByUserID(userID, page, limit)
}

func (s *orderService) CancelOrder(userID, orderID uint) error {
	// Get order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	// Verify ownership
	if order.UserID != userID {
		return errors.New("order not found")
	}

	// Check if order can be cancelled
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusProcessing {
		return errors.New("order cannot be cancelled")
	}

	// Use transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Restore stock
		for _, item := range order.Items {
			product, err := s.productRepo.FindByID(item.ProductID)
			if err != nil {
				continue // Product might be deleted
			}

			if product.TrackInventory {
				if err := s.productRepo.IncrementStock(item.ProductID, item.Quantity); err != nil {
					return errors.New("failed to restore stock")
				}
			}
		}

		// Update order status
		return s.orderRepo.UpdateStatus(orderID, domain.OrderStatusCancelled)
	})
}

func (s *orderService) GetAllOrders(query *domain.OrderListQuery) ([]*domain.Order, int64, error) {
	return s.orderRepo.List(query)
}

func (s *orderService) UpdateOrderStatus(orderID uint, status domain.OrderStatus) error {
	// Verify order exists
	_, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	return s.orderRepo.UpdateStatus(orderID, status)
}
