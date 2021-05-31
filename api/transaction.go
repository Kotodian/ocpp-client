package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ocpp-client/websocket"
)

// PlugCable 插入电缆
func PlugCable(c *gin.Context) {
	request := &struct {
		SN string `json:"sn"`
	}{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(request.SN) == 0 {
		c.JSON(http.StatusBadRequest, errors.New("invalid sn"))
		return
	}
	// 确认连接是否存在
	_ws, ok := websocket.Cache.Get(request.SN)
	if !ok {
		c.JSON(http.StatusBadRequest, errors.New("this charge station doesn't exist"))
		return
	}
	// 断言
	ws := _ws.(*websocket.Client)
	// 如果该连接未连接到平台
	if !ws.Connected() {
		c.JSON(http.StatusBadRequest, errors.New("this charge station is offline"))
		return
	}
	err = ws.Instance().StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
