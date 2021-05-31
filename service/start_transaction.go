package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) RequestStartTransactionResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.RequestStartTransactionRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.RequestStartTransactionResponseJson{
		Status: message.RequestStartStopStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "RemoteStartTransaction", response, msgID)
	return msg, err
}
