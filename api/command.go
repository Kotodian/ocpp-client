package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ocpp-client/message"
	"ocpp-client/websocket"
)

// FirmwareStatusNotification 固件状态上报
// @method: post
// @group: Command
// @router: /firmware_status_notification/:sn
func FirmwareStatusNotification(c *gin.Context) {
	// 绑定request json
	request := &message.FirmwareStatusNotificationRequestJson{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	conn := paramToConn(c)
	if conn == nil {
		c.JSON(http.StatusBadRequest, "conn doesn't exist")
		return
	}
	sendMessage(conn, "FirmwareStatusNotification", request)
	c.JSON(http.StatusOK, "success")
}

// LogStatusNotification 日志状态上报
// @method: post
// @group: Command
// @router: /log_status_notification/:sn
func LogStatusNotification(c *gin.Context) {
	request := &message.LogStatusNotificationRequestJson{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	conn := paramToConn(c)
	if conn == nil {
		c.JSON(http.StatusBadRequest, "conn doesn't exist")
		return
	}
	sendMessage(conn, "LogStatusNotification", request)
	c.JSON(http.StatusOK, "success")
}

// NotifyReport 上报基本数据
// @method: post
// @group: Command
// @router: /notify_report/:sn
func NotifyReport(c *gin.Context) {
	request := &message.NotifyReportRequestJson{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	conn := paramToConn(c)
	if conn == nil {
		c.JSON(http.StatusBadRequest, "conn doesn't exist")
		return
	}
	sendMessage(conn, "NotifyReport", request)
	c.JSON(http.StatusOK, "success")
}
func paramToConn(c *gin.Context) *websocket.Client {
	// 获取路由参数sn
	sn := c.Param("sn")
	// 获取websocket连接
	_chargeStation, exists := websocket.Get(sn)
	if !exists {
		return nil
	}
	chargeStation, ok := _chargeStation.(*websocket.Client)
	if !ok {
		return nil
	}
	return chargeStation
}

func sendMessage(conn *websocket.Client, action string, payload interface{}) {
	msg, _, _ := message.New("2", action, payload)
	conn.Write(msg)
}
