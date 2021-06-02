package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) RequestStopTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	c.lock.Lock()
	request := &message.RequestStopTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		c.StopTransaction()
		_ = DB.Put(ChargeStationBucket, c.ID(), c)
		c.lock.Unlock()
		time.Sleep(1 * time.Second)
	}()
	response := &message.RequestStopTransactionResponseJson{
		Status: message.RequestStartStopStatusEnumType_3_Accepted,
	}
	msg, _, _ = message.New("3", "RequestStopTransaction", response, msgID)
	return msg, nil
}

// StopTransaction 关闭充电
func (c *ChargeStation) StopTransaction() {
	if c.Transaction == nil ||
		c.Transaction.EventType == message.TransactionEventEnumType_1_Ended {
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
	c.Transaction.Instance.StoppedReason = &reason
	c.Transaction.EventType = message.TransactionEventEnumType_1_Ended
	c.Electricity = minElectricity
	c.SendEvent()
}
