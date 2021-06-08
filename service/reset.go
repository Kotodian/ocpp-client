package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) ResetResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.ResetRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.ResetResponseJson{
		Status: message.ResetStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "Reset", response, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}
