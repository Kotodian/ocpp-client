package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) RequestStartTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	c.lock.Lock()
	request := &message.RequestStartTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		c.lock.Unlock()
		if err != nil {
			return
		}
		// 转换成updated并定时发送
		if c.Transaction.EventType == message.TransactionEventEnumType_1_Started {
			_ = c.StartTransaction()
		}
	}()

	response := &message.RequestStartTransactionResponseJson{}
	if c.Transaction.EventType == message.TransactionEventEnumType_1_Updated ||
		c.Connectors[0].State == message.ConnectorStatusEnumType_1_Occupied {
		response.Status = message.RequestStartStopStatusEnumType_1_Rejected
		goto send
	} else if c.Transaction.EventType == message.TransactionEventEnumType_1_Started {
		response.Status = message.RequestStartStopStatusEnumType_1_Accepted
		response.TransactionId = &c.Transaction.Instance.TransactionId
		goto send
	} else if c.Transaction == nil ||
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
		err = DB.Put(ChargeStationBucket, c.ID(), c)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
send:
	msg, _, _ = message.New("3", "RemoteStartTransaction", response, msgID)
	c.Resend <- msg
	return nil, nil
}
