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
			_ = c.StartTransaction()
			time.Sleep(1 * time.Second)
		}
		c.transaction.instance.RemoteStartId = &request.RemoteStartId
		c.transaction.idTokenType = message.IdTokenEnumType_7_Central
		c.transaction.idToken = &request.IdToken
		_ = c.StartTransaction()
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
