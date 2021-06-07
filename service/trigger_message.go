package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) TriggerMessageResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.TriggerMessageRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.TriggerMessageResponseJson{
		Status: message.TriggerMessageStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "TriggerMessage", response, msgID)
	return msg, err
}
