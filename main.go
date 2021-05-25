package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"ocpp-client/api"
	"os"
)

func main() {
	engine := gin.Default()
	engine.POST("/create", api.NewChargeStation)
	engine.POST("/command", api.Command)
	_ = engine.Run(":" + os.Getenv("PORT"))
}
