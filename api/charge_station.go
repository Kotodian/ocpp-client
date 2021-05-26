package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/url"
	"ocpp-client/service"
	"ocpp-client/websocket"
	"os"
	"strings"
	"time"
)

func NewChargeStation(c *gin.Context) {
	request := &struct {
		// 前缀
		SN string `json:"sn"`
		// 个数
		Nums int `json:"nums"`
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
	if len(request.SN) == 0 && !strings.HasPrefix(request.SN, "T") {
		c.JSON(400, errors.New("invalid sn or addr"))
		return
	}
	// 默认启动一个
	if request.Nums <= 0 {
		request.Nums = 1
	}
	var g errgroup.Group
	for i := 1; i <= request.Nums; i++ {
		g.Go(func() error {
			sn := request.SN + fmt.Sprint(i)
			sn = formatSN(sn)
			station := service.NewChargeStation(sn)
			client := websocket.NewClient(station)
			addr := url.URL{Scheme: "ws", Host: os.Getenv("ADDR"), Path: "/ocpp/" + sn}
			err := client.Conn(addr.String())
			if err != nil {
				return fmt.Errorf("%s connect error %s", sn, err.Error())
			}
			return nil
		})

	}
	if request.Sleep <= 0 {
		time.Sleep(100 * time.Millisecond)
	} else {
		time.Sleep(time.Duration(request.Sleep) * time.Millisecond)
	}
	if err = g.Wait(); err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "success")
}

func Command(c *gin.Context) {
	request := &struct {
		SN      string `json:"sn"`
		Command string `json:"command"`
	}{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if len(request.SN) == 0 || len(request.Command) == 0 {
		c.JSON(400, errors.New("invalid sn or command"))
		return
	}
	_websocket, ok := websocket.Cache.Get(request.SN)
	if !ok {
		c.JSON(400, errors.New("this charge station doesn't exist"))
		return
	}
	// 断言
	ws := _websocket.(*websocket.Client)
	if ws.Connected() {
		msg, err := ws.Instance().Function("2", "", request.Command)
		if err != nil {
			c.JSON(500, err)
			return
		}
		ws.Write(msg)
		c.JSON(200, "success")
		return
	} else {
		c.JSON(500, errors.New("this charge station doesn't connect"))
		return
	}
}

func formatSN(sn string) string {
	if len(sn) > 11 {
		sn = sn[:11]
	} else if len(sn) < 11 {
		repeat := strings.Repeat("0", 11-len(sn))
		sn += repeat
	}
	return sn
}
