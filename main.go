package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"ocpp-client/api"
	"os"
)

func main() {
	engine := gin.Default()
	// 创建charge station websocket
	// ps. sn是填前缀后续全部自己生成
	engine.POST("/create", api.NewChargeStation)
	// 由桩主动发送命令
	engine.POST("/command", api.Command)
	// 插进电缆生成transaction
	engine.POST("/plugin", api.PlugCable)
	_ = engine.Run(":" + os.Getenv("PORT"))
}
