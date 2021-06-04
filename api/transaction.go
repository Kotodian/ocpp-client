package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ocpp-client/websocket"
)

// TransactionEvent 充电
// @method: post
// @group: /transaction
// @router: /add
func TransactionEvent(c *gin.Context) {
	request := &struct {
		SN string `json:"sn"`
	}{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	_conn, ok := websocket.Cache.Get(request.SN)
	if !ok {
		c.JSON(http.StatusBadRequest, errors.New("no conn to ac-ocpp"))
		return
	}
	conn := _conn.(*websocket.Client)
	// 如果未连接
	if !conn.Connected() {
		c.JSON(http.StatusBadRequest, errors.New("no conn to ac-ocpp"))
		return
	}
	conn.Instance().Lock()
	var inTransaction bool
	defer func() {
		conn.Instance().Unlock()
		if !inTransaction {
			_ = conn.Instance().UpdateTransaction()
		}
	}()
	// 如果在充电中
	if conn.Instance().InTransaction() {
		inTransaction = true
		c.JSON(http.StatusBadRequest, errors.New("in transaction"))
		return
	}

	_ = conn.Instance().StartTransaction()
}
