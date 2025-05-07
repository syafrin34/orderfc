package main

import (
	"fmt"
	"orderfc/cmd/order/resource"
	"orderfc/config"
	"orderfc/infrastructure/logger"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
	//db := resource.InitDB(&cfg)
	//redis := resource.InitRedis(&cfg)
	resource.InitDB(&cfg)
	resource.InitRedis(&cfg)
	logger.SetupLogger()

}
