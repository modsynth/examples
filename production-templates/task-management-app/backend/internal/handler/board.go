package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-management-app/internal/domain"
	"task-management-app/internal/service"
)

type BoardHandler struct {
	boardService service.BoardService
}

func NewBoardHandler(boardService service.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	projectID, err := strconv.ParseUint(c.Param("projectID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	var req domain.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board, err := h.boardService.Create(uint(projectID), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func (h *BoardHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	board, err := h.boardService.GetByID(uint(boardID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	var req domain.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board, err := h.boardService.Update(uint(boardID), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	if err := h.boardService.Delete(uint(boardID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "board deleted successfully"})
}

func (h *BoardHandler) ListByProject(c *gin.Context) {
	userID := c.GetUint("userID")
	projectID, err := strconv.ParseUint(c.Param("projectID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	boards, err := h.boardService.ListByProject(uint(projectID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boards)
}
