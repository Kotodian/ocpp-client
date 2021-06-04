// Package init Code generated by "router-annotation";DO NOT EDIT.
package init

import (
	"github.com/gin-gonic/gin"
	"ocpp-client/api"
)

var engine *gin.Engine

func init() {
	engine = gin.Default()
	stationGroup := engine.Group("/station")
	{
		stationGroup.POST("/list", api.ListChargeStation)
	}
	{
		stationGroup.POST("/list", api.NewChargeStation)
	}
	{
		stationGroup.POST("/command", api.Command)
	}
	transactionGroup := engine.Group("/transaction")
	{
		transactionGroup.POST("/add", api.TransactionEvent)
	}
}
func GinEngine() *gin.Engine {
	return engine
}
