package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/service"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart godoc
// @Summary Get user's cart
// @Tags cart
// @Produce json
// @Success 200 {object} domain.CartWithSummary
// @Failure 401 {object} map[string]string
// @Router /api/v1/cart [get]
// @Security BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	cart, err := h.cartService.GetCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart godoc
// @Summary Add item to cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body domain.AddToCartRequest true "Cart item"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/cart/items [post]
// @Security BearerAuth
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.AddToCart(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Tags cart
// @Accept json
// @Produce json
// @Param id path int true "Cart Item ID"
// @Param request body domain.UpdateCartItemRequest true "Quantity"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/cart/items/{id} [put]
// @Security BearerAuth
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, _ := c.Get("user_id")

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	var req domain.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.UpdateCartItem(userID.(uint), uint(itemID), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart item updated"})
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Tags cart
// @Param id path int true "Cart Item ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/v1/cart/items/{id} [delete]
// @Security BearerAuth
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := h.cartService.RemoveFromCart(userID.(uint), uint(itemID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ClearCart godoc
// @Summary Clear user's cart
// @Tags cart
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/v1/cart [delete]
// @Security BearerAuth
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.cartService.ClearCart(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
