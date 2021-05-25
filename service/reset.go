package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func ResetResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.ResetRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.ResetResponseJson{
		Status: message.ResetStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "Reset", response, msgID)
	return msg, err
}
