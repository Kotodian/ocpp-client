package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"ocpp-client/api"
	"os"
)

func main() {
	engine := gin.Default()
	// 充电桩列表
	engine.POST("/station", api.ListChargeStation)
	// 创建charge station websocket
	// ps. sn是填前缀后续全部自己生成
	engine.POST("/create", api.NewChargeStation)
	// 由桩主动发送命令
	engine.POST("/command", api.Command)
	// 桩主动去充电
	engine.POST("/transaction", api.TransactionEvent)
	_ = engine.Run(":" + os.Getenv("PORT"))
}
