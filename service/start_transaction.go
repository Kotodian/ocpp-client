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
	defer c.lock.Unlock()

	response := &message.RequestStartTransactionResponseJson{}
	if c.transaction.eventType == message.TransactionEventEnumType_1_Started ||
		c.transaction.eventType == message.TransactionEventEnumType_1_Updated ||
		c.connectors[0].State() == message.ConnectorStatusEnumType_1_Occupied {
		response.Status = message.RequestStartStopStatusEnumType_1_Rejected
		goto send
	}
	if c.transaction == nil ||
		c.transaction.eventType == message.TransactionEventEnumType_1_Ended {
		// 创建充电事件
		_ = c.StartTransaction()
		response.Status = message.RequestStartStopStatusEnumType_1_Accepted
		response.TransactionId = &c.transaction.instance.TransactionId
		// 先发送RequestStartTransactionResponse
		msg, _, _ = message.New("3", "RemoteStartTransaction", response, msgID)
		c.Resend <- msg
		time.Sleep(1 * time.Second)
		// 再发送TransactionEventRequest
		c.transaction.instance.RemoteStartId = &request.RemoteStartId
		c.transaction.idTokenType = message.IdTokenEnumType_7_Central
		c.transaction.idToken = &request.IdToken
		c.SendEvent()
		time.Sleep(1 * time.Second)
		// 然后开始定时发送TransactionEventRequest
		_ = c.StartTransaction()
	}
send:
	msg, _, _ = message.New("3", "RemoteStartTransaction", response, msgID)
	c.Resend <- msg
	return nil, nil
}
