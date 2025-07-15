package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"internal-transfer-system/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req service.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.CreateAccount(&req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "account already exists"):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "positive integer") ||
			strings.Contains(err.Error(), "non-negative") ||
			strings.Contains(err.Error(), "decimal places") ||
			strings.Contains(err.Error(), "invalid"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.service.GetAccount(uint(accountID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req service.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.CreateTransaction(&req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "insufficient balance"):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "same account") ||
			strings.Contains(err.Error(), "greater than 0") ||
			strings.Contains(err.Error(), "decimal places") ||
			strings.Contains(err.Error(), "invalid"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}