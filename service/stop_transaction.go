package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) RequestStopTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.RequestStopTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		time.Sleep(1 * time.Second)
		c.StopTransaction()
	}()
	response := &message.RequestStopTransactionResponseJson{
		Status: message.RequestStartStopStatusEnumType_3_Accepted,
	}
	msg, _, err = message.New("3", "RequestStopTransaction", response, msgID)
	return msg, err
}
