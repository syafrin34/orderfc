package routes

import (
	"orderfc/cmd/order/handler"
	"orderfc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, orderHandler handler.OrderHandler, jwtsecret string) {
	router.Use(middleware.RequestLogger())
	authMiddleware := middleware.AuthMiddleware(jwtsecret)
	router.Use(authMiddleware)
	router.POST("/v1/checkout", orderHandler.CheckOutOrder)
	router.GET("/v1/order_history", orderHandler.GetOrderHistory)
}
