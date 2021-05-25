package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) FirmwareStatusNotificationRequest() ([]byte, error) {
	requestId := 1
	request := message.FirmwareStatusNotificationRequestJson{
		RequestId: &requestId,
		Status:    message.FirmwareStatusEnumType_1_Idle,
	}
	msg, _, err := message.New("2", "FirmwareStatusNotification", request)
	return msg, err
}

func (c *ChargeStation) FirmwareStatusNotificationResponse(msgID string, msg []byte) error {
	response := &message.FirmwareStatusNotificationResponseJson{}
	return json.Unmarshal(msg, response)
}
