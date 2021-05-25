package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) RequestStopTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.RequestStopTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.RequestStopTransactionResponseJson{
		Status: message.RequestStartStopStatusEnumType_3_Accepted,
	}
	msg, _, err = message.New("3", "RequestStopTransaction", response, msgID)
	return msg, err
}
