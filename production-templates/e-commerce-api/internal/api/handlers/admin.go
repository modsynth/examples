package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/service"
)

type AdminHandler struct {
	orderService service.OrderService
}

func NewAdminHandler(orderService service.OrderService) *AdminHandler {
	return &AdminHandler{
		orderService: orderService,
	}
}

// GetAllOrders godoc
// @Summary Get all orders (Admin only)
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Filter by status"
// @Success 200 {array} domain.Order
// @Router /api/v1/admin/orders [get]
// @Security BearerAuth
func (h *AdminHandler) GetAllOrders(c *gin.Context) {
	var query domain.OrderListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, total, err := h.orderService.GetAllOrders(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body domain.UpdateOrderStatusRequest true "Order status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/orders/{id} [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req domain.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderStatus(uint(orderID), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

// GetStats godoc
// @Summary Get dashboard statistics (Admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/stats [get]
// @Security BearerAuth
func (h *AdminHandler) GetStats(c *gin.Context) {
	// Placeholder for dashboard statistics
	// This would typically aggregate data from multiple sources
	stats := gin.H{
		"total_orders":    0,
		"total_revenue":   0.0,
		"pending_orders":  0,
		"total_customers": 0,
		"total_products":  0,
	}

	c.JSON(http.StatusOK, stats)
}
