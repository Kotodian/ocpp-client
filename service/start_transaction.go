package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) RequestStartTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.RequestStartTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		time.Sleep(1 * time.Second)
		if c.transaction == nil {
			_ = c.StartTransaction(request.RemoteStartId)
			time.Sleep(1 * time.Second)
		}
		_ = c.StartTransaction(request.RemoteStartId)
	}()

	response := &message.RequestStartTransactionResponseJson{
		Status: message.RequestStartStopStatusEnumType_1_Accepted,
	}
	if c.transaction != nil {
		response.TransactionId = &c.transaction.instance.TransactionId
	}
	msg, _, err = message.New("3", "RemoteStartTransaction", response, msgID)
	return msg, err
}
