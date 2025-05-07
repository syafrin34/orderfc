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

}
