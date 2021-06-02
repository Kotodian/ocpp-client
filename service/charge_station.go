package service

import (
	"errors"
	"ocpp-client/message"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// ChargeStation 充电桩实例 具有一些属性
type ChargeStation struct {
	// 唯一值
	Sn string `json:"sn"`
	// 运营商名称
	VendorName string `json:"vendor_name"`
	// 不知道干啥用的
	Model string `json:"model"`
	// heartbeat时间
	interval time.Duration
	// 停止channel
	stop chan struct{}
	// 重新发送命令的channel
	Resend chan []byte `json:"-"`
	// 锁
	lock sync.Mutex
	// 充电枪
	Connectors []*Connector `json:"connectors"`
	// 正在执行的transaction
	Transaction *Transaction `json:"transaction"`
	// 电量
	Electricity float64 `json:"electricity"`
}

// NewChargeStation 通过sn创建实例
func NewChargeStation(sn string) *ChargeStation {
	chargeStation := &ChargeStation{
		Sn:          sn,
		stop:        make(chan struct{}),
		VendorName:  "JoysonQuin",
		Model:       "JWBOX",
		Resend:      make(chan []byte, 1),
		Connectors:  make([]*Connector, 0),
		Electricity: minElectricity,
	}
	chargeStation.Connectors = append(chargeStation.Connectors, NewConnector(1))
	_ = DB.Put(ChargeStationBucket, sn, chargeStation)
	// 建立一个默认的充电枪
	return chargeStation
}

// ID 充电桩唯一值
func (c *ChargeStation) ID() string {
	return c.Sn
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
		//定时发送heartbeat命令
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
		time.Sleep(100 * time.Millisecond)
		if c.Transaction == nil {
			return
		}
		if c.Transaction.eventType == message.TransactionEventEnumType_1_Updated ||
			c.Transaction.eventType == message.TransactionEventEnumType_1_Started {
			msg, _ := c.TransactionEventRequest()
			c.Resend <- msg
		}
	}
	c.stop = make(chan struct{})
}

// StartTransaction 开始充电
func (c *ChargeStation) StartTransaction() error {
	if c.Transaction != nil && c.Transaction.eventType != message.TransactionEventEnumType_1_Ended {
		if c.Transaction.eventType == message.TransactionEventEnumType_1_Started {
			state := message.ChargingStateEnumType_1_Charging
			c.Transaction.eventType = message.TransactionEventEnumType_1_Updated
			c.Transaction.instance.ChargingState = &state
			_, _ = c.TransactionEventRequest()
		} else {
			return errors.New("in transaction")
		}
		return nil
	} else if c.Connectors[0].State() == message.ConnectorStatusEnumType_1_Available {
		c.Connectors[0].SetState(message.ConnectorStatusEnumType_1_Occupied)
		// 通知平台枪的状态发生改变
		msg, _ := c.StatusNotificationRequest()
		// 发送给平台
		c.Resend <- msg
		// 等待一段时间接收response
		time.Sleep(100 * time.Millisecond)
		// 新建一个id
		id := strconv.FormatInt(time.Now().Unix(), 10)
		instance := &message.TransactionType{
			TransactionId: id,
		}
		transaction := NewTransaction(instance)
		c.Transaction = transaction
		return nil
	} else {
		return errors.New("in transaction")
	}
}

// StopTransaction 关闭充电
func (c *ChargeStation) StopTransaction() {
	if c.Transaction == nil ||
		c.Transaction.eventType == message.TransactionEventEnumType_1_Ended {
		return
	}
	// 阻塞到直到transactionEvent接收到值
	c.Transaction.stop <- struct{}{}
	// 拔枪
	time.Sleep(100 * time.Millisecond)
	c.Connectors[0].SetState(message.ConnectorStatusEnumType_1_Available)
	msg, _ := c.StatusNotificationRequest()
	c.Resend <- msg
	// 发送关闭
	// unplug cable
	time.Sleep(100 * time.Millisecond)
	reason := message.ReasonEnumType_1_Remote
	c.Transaction.instance.StoppedReason = &reason
	c.Transaction.eventType = message.TransactionEventEnumType_1_Ended
	c.Electricity = minElectricity
	c.SendEvent()
}
