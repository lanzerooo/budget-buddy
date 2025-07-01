package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/model"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes(r *gin.Engine) {
	r.POST("/api/transactions", h.createTransaction)
	r.GET("/api/transactions", h.getTransactions)
	r.GET("/api/balance", h.getBalance)
}

func (h *Handler) createTransaction(c *gin.Context) {
	var tx model.Transaction
	if err := c.BindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	tx.CreatedAt = time.Now()
	if err := h.service.CreateTransaction(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (h *Handler) getTransactions(c *gin.Context) {
	list, err := h.service.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) getBalance(c *gin.Context) {
	balance, err := h.service.GetBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
