package api

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"ocpp-client/service"
	"ocpp-client/websocket"
	"os"
	"strings"
	"time"
)

// ListChargeStation 充电桩列表
// @method: post
// @group: ChargeStation
// @router: /list
func ListChargeStation(c *gin.Context) {
	list, err := service.ListChargeStation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// NewChargeStation 创建充电桩
// @method: post
// @group: ChargeStation
// @router: /add
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
	go func() {
		for i := 1; i <= request.Nums; i++ {
			var sn string
			if len(request.SN) != 11 {
				sn = request.SN + createRandomNumber(11-len(request.SN))
				sn = formatSN(sn)
			} else {
				sn = request.SN
			}
			fmt.Println(request.SN)
			station := service.NewChargeStation(sn)
			//exists, err := service.DB.Exists(station.BucketName(), station.ID())
			//if err != nil {
			//	log.Println(err)
			//	return
			//}
			//if !exists {
			//	err = service.DB.Put(sn, station)
			//}
			//if err != nil {
			//	log.Println(err)
			//	return
			//}
			client := websocket.NewClient(station)
			if client == nil {
				continue
			}
			addr := url.URL{Scheme: "ws", Host: os.Getenv("ADDR"), Path: "/ocpp/" + sn}
			//if exists {
			//	err = client.ReConn()
			//} else {
			err = client.Conn(addr.String())
			//}
			if err != nil {
				log.Printf("%s connect error %s\n", sn, err.Error())
				continue
			}

			if request.Sleep <= 0 {
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(time.Duration(request.Sleep) * time.Millisecond)
			}
		}
	}()

	c.JSON(200, "success")
}

// Command 发送充电桩命令
// @method: post
// @group: ChargeStation
// @router: /command
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

func createRandomNumber(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

// CloseChargeStation 关闭websocket连接
// @method: post
// @group: ChargeStation
// @router: /close
func CloseChargeStation(c *gin.Context) {
	request := &struct {
		SN string `json:"sn"`
	}{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(request.SN) == 0 {
		c.JSON(http.StatusBadRequest, errors.New("no sn"))
		return
	}
	_conn, ok := websocket.Get(request.SN)
	if !ok {
		c.JSON(http.StatusBadRequest, errors.New("no conn"))
		return
	}
	conn, ok := _conn.(*websocket.Client)
	if !ok {
		c.JSON(http.StatusBadRequest, errors.New("not conn"))
		return
	}
	conn.Close()
	c.JSON(http.StatusOK, "success")
}
