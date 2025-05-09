package main

import (
	"fmt"
	"orderfc/cmd/order/handler"
	"orderfc/cmd/order/repository"
	"orderfc/cmd/order/resource"
	"orderfc/cmd/order/service"
	"orderfc/cmd/order/usecase"
	"orderfc/config"
	"orderfc/infrastructure/logger"
	"orderfc/kafka"
	"orderfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	// resource.InitDB(&cfg)
	// resource.InitRedis(&cfg)
	logger.SetupLogger()
	kafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9093"}, "order.created")
	defer kafkaProducer.Close()

	orderRepository := repository.NewOrderRepository(db, redis, cfg.Product.Host)
	orderService := service.NewOrderService(*orderRepository)
	orderUsecase := usecase.NewOrderUsecase(*orderService, *kafkaProducer)
	orderhandler := handler.NewOrderHandler(*orderUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *orderhandler, cfg.Secret.JWTSecret)
	router.Run(":" + port)

}
