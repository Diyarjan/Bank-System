package handler

import (
	"github.com/Diyarjan/BankSystem/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	control := router.Group("api/accounts")
	{
		control.POST("/", h.createAccount)
		control.DELETE("/:account_id", h.deleteAccount)
		control.GET("/:account_id", h.getAccountByID)
		control.GET("/", h.getAllAccounts)
	}

	transaction := router.Group("/api/accounts/:account_id")
	{
		transaction.POST("/deposit", h.depositToAccount)
		transaction.POST("/withdraw", h.withdrawFromAccount)
		transaction.POST("/transfer", h.transferFunds)
		transaction.GET("/transactions", h.getTransactionHistory)
	}

	return router
}
