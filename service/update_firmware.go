package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) UpdateFirmwareResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.UpdateFirmwareRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.UpdateFirmwareResponseJson{
		Status: message.UpdateFirmwareStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "UpdateFirmware", response, msgID)
	return msg, err
}
