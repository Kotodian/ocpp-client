package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"ocpp-client/service"
	"ocpp-client/websocket"
	"time"
)

func NewChargeStation(c *gin.Context) {
	request := struct {
		// 前缀
		SN string `json:"sn"`
		// 个数
		Nums int `json:"nums"`
		// 连接地址
		Addr string `json:"addr"`
		// 睡眠时间(ms)
		Sleep int `json:"sleep"`
	}{}
	// 绑定参数
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(400, err)
		return
	}
	// 校验
	if len(request.SN) == 0 || len(request.Addr) == 0 || request.Nums <= 0 {
		c.JSON(400, errors.New("invalid sn or nums"))
		return
	}
	var g errgroup.Group
	for i := 1; i <= request.Nums; i++ {
		g.Go(func() error {
			sn := request.SN + fmt.Sprint(i)
			station := service.NewChargeStation(sn)
			client := websocket.NewClient(station)
			err := client.Conn("ws://" + request.Addr + "/ocpp/" + sn)
			if err != nil {
				return fmt.Errorf("%s connect error %s", sn, err.Error())
			}
			return nil
		})
		time.Sleep(time.Duration(request.Sleep) * time.Millisecond)
	}
	if err = g.Wait(); err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "success")
}
