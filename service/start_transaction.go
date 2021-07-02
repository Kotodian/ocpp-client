package service

import (
	"encoding/json"
	"errors"
	"ocpp-client/message"
	"strconv"
	"time"
)

func (c *ChargeStation) RequestStartTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	request := &message.RequestStartTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	response := &message.RequestStartTransactionResponseJson{}
	defer func() {
		if err != nil {
			return
		}
		if response.Status == message.RequestStartStopStatusEnumType_1_Accepted {
			// 转换成updated并定时发送
			if c.Transaction.EventType == message.TransactionEventEnumType_1_Started {
				_ = c.UpdateTransaction()
			}
		}
	}()

	if c.InTransaction() {
		response.Status = message.RequestStartStopStatusEnumType_1_Rejected
		goto send
	}

	if c.Transaction == nil ||
		c.Transaction.EventType == message.TransactionEventEnumType_1_Ended {
		// 创建充电事件
		_ = c.StartTransaction()
		response.Status = message.RequestStartStopStatusEnumType_1_Accepted
		response.TransactionId = &c.Transaction.Instance.TransactionId
		// 先发送RequestStartTransactionResponse
		msg, _, _ = message.New("3", "RemoteStartTransaction", response, msgID)
		c.Resend <- msg
		time.Sleep(1 * time.Second)
		// 再发送TransactionEventRequest
		c.Transaction.Instance.RemoteStartId = &request.RemoteStartId
		c.Transaction.IdTokenType = message.IdTokenEnumType_7_Central
		c.Transaction.IdToken = &request.IdToken
		c.SendEvent()
		// 存进数据库
		//err = DB.Put(c.ID(), c)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else if c.Transaction.EventType == message.TransactionEventEnumType_1_Updated ||
		c.Connectors[0].State == message.ConnectorStatusEnumType_1_Occupied {
		response.Status = message.RequestStartStopStatusEnumType_1_Rejected
		goto send
	} else if c.Transaction.EventType == message.TransactionEventEnumType_1_Started {
		response.Status = message.RequestStartStopStatusEnumType_1_Accepted
		response.TransactionId = &c.Transaction.Instance.TransactionId
		goto send
	}
send:
	msg, _, _ = message.New("3", "RemoteStartTransaction", response, msgID)
	c.Resend <- msg
	return nil, nil
}

// StartTransaction 开始充电
func (c *ChargeStation) StartTransaction() error {
	if c.Connectors[0].State == message.ConnectorStatusEnumType_1_Available {
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

func (c *ChargeStation) UpdateTransaction() error {
	if c.Transaction != nil && c.Transaction.EventType != message.TransactionEventEnumType_1_Ended {
		if c.Transaction.EventType == message.TransactionEventEnumType_1_Started {
			state := message.ChargingStateEnumType_1_Charging
			c.Transaction.EventType = message.TransactionEventEnumType_1_Updated
			c.Transaction.Instance.ChargingState = &state
			_, _ = c.TransactionEventRequest()
		} else {
			return errors.New("in transaction")
		}
	}
	return nil
}
