package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) GetLogResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.GetLogRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.GetLogResponseJson{
		Status: message.LogStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "GetLog", response, msgID)
	return msg, err
}
