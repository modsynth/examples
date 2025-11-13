package domain

import "time"

type OrderStatus string
type PaymentStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Order struct {
	ID                     uint          `json:"id" gorm:"primaryKey"`
	UserID                 uint          `json:"user_id" gorm:"not null"`
	User                   *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderNumber            string        `json:"order_number" gorm:"uniqueIndex;not null"`
	Status                 OrderStatus   `json:"status" gorm:"not null;default:'pending'"`
	Subtotal               float64       `json:"subtotal" gorm:"not null"`
	Tax                    float64       `json:"tax" gorm:"not null;default:0"`
	Shipping               float64       `json:"shipping" gorm:"not null;default:0"`
	Total                  float64       `json:"total" gorm:"not null"`
	Currency               string        `json:"currency" gorm:"not null;default:'USD'"`
	PaymentStatus          PaymentStatus `json:"payment_status" gorm:"not null;default:'pending'"`
	PaymentMethod          string        `json:"payment_method"`
	StripePaymentIntentID  string        `json:"stripe_payment_intent_id"`
	ShippingAddressLine1   string        `json:"shipping_address_line1"`
	ShippingAddressLine2   string        `json:"shipping_address_line2"`
	ShippingCity           string        `json:"shipping_city"`
	ShippingState          string        `json:"shipping_state"`
	ShippingPostalCode     string        `json:"shipping_postal_code"`
	ShippingCountry        string        `json:"shipping_country"`
	Notes                  string        `json:"notes"`
	Items                  []OrderItem   `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
}

type OrderItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrderID     uint      `json:"order_id" gorm:"not null"`
	ProductID   uint      `json:"product_id" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"not null"`
	ProductSKU  string    `json:"product_sku"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Subtotal    float64   `json:"subtotal" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShippingAddress struct {
	Line1      string `json:"line1" binding:"required"`
	Line2      string `json:"line2"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required,len=2"`
}

type CreateOrderRequest struct {
	ShippingAddress ShippingAddress `json:"shipping_address" binding:"required"`
	PaymentMethod   string          `json:"payment_method" binding:"required"`
	Notes           string          `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}

type OrderListQuery struct {
	Page          int            `form:"page" binding:"omitempty,gte=1"`
	Limit         int            `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Status        *OrderStatus   `form:"status"`
	PaymentStatus *PaymentStatus `form:"payment_status"`
	UserID        *uint          `form:"user_id"`
	OrderNumber   string         `form:"order_number"`
}
