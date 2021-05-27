package service

import (
	"reflect"
	"sync"
	"time"
)

// ChargeStation 充电桩实例 具有一些属性
type ChargeStation struct {
	// 唯一值
	sn string
	// 运营商名称
	vendorName string
	// 不知道干啥用的
	model string
	// heartbeat时间
	interval time.Duration
	// 停止channel
	stop chan struct{}
	// 重新发送命令的channel
	Resend chan []byte
	// 充电事务通道
	transactions chan Transaction
	// 锁
	lock sync.Mutex
	// 目前处理的事务
	transaction *Transaction
}

// NewChargeStation 通过sn创建实例
func NewChargeStation(sn string) *ChargeStation {
	chargeStation := &ChargeStation{
		sn:           sn,
		stop:         make(chan struct{}, 1),
		vendorName:   "joysonquin",
		model:        "test",
		Resend:       make(chan []byte, 1),
		transactions: make(chan Transaction, 1),
	}
	return chargeStation
}

// ID 充电桩唯一值
func (c *ChargeStation) ID() string {
	return c.sn
}

// Function 调用指定命令 parameters代表参数一般就是payload,由桩主动发起的命令不需要参数以及messageID
// example c.Function("2","","BootNotification") 或者 c.Function("3", "BootNotification",payload, messageID)
func (c *ChargeStation) Function(typ string, messageID, action string, parameters ...interface{}) ([]byte, error) {
	// 确认消息类型
	if typ == "2" {
		action += "Request"
	} else {
		action += "Response"
	}
	// 参数列表
	var p []reflect.Value
	// 如果需要提供messageID,就将messageID放入到参数列表中
	if len(messageID) > 0 {
		p = append(p, reflect.ValueOf(messageID))
	}
	// 如果有其他参数,就加入到参数列表中
	if len(parameters) > 0 {
		for _, parameter := range parameters {
			p = append(p, reflect.ValueOf(parameter))
		}
	}
	//根据charge station的反射调用方法
	// values 表示返回结果 如果该返回值有两个,第一个就是[]byte类型,第二个代表error
	// 如果返回值只有一个就表示只有error
	values := reflect.ValueOf(c).MethodByName(action).Call(p)
	if len(values) == 2 {
		// 判断error为空,形如if err == nil {return value, nil}
		if values[1].IsNil() {
			return values[0].Bytes(), nil
		} else {
			// 断言 error
			return nil, values[1].Interface().(error)
		}

	} else {
		if values[0].IsNil() {
			return nil, nil
		} else {
			return nil, values[0].Interface().(error)
		}
	}
}

// Stop 停止信号 如果该channel接收到值就需要停止发送心跳数据
func (c *ChargeStation) Stop() {
	close(c.stop)
}

// ReConn 重新连接
func (c *ChargeStation) ReConn() {
	// 如果interval等于0表示该charge station还未接收到ac-ocpp的BootNotification的response就断开连接了
	// 所以就重新发送BootNotificationRequest
	if c.interval == 0 {
		msg, _ := c.BootNotificationRequest()
		c.Resend <- msg
	} else {
		// 定时发送heartbeat命令
		go func() {
			ticker := time.NewTicker(c.interval)
			defer ticker.Stop()
			for {
				select {
				// 如果停止了就关闭Heartbeat
				case <-c.stop:
					return
				// 时间到了就发送Heartbeat
				case <-ticker.C:
					msg, _ := c.HeartbeatRequest()
					c.Resend <- msg
				}
			}
		}()
		// 发送StatusNotification
		msg, _ := c.StatusNotificationRequest()
		c.Resend <- msg
	}
}

func (c *ChargeStation) SendTransactionEvent() {
	for {
		select {
		case <-c.stop:
			return
		case transaction := <-c.transactions:
			c.lock.Lock()
			c.transaction = &transaction
			event, err := c.TransactionEventRequest()
			if err != nil {
				goto unlock
			}
			if event != nil {
				c.Resend <- event
			}
		unlock:
			c.lock.Unlock()
		}
	}
}
