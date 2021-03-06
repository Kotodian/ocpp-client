package service

import (
	"encoding/json"
	"ocpp-client/message"
)

// LogStatusNotificationRequest log status
func (c *ChargeStation) LogStatusNotificationRequest() ([]byte, error) {
	id := 9
	request := &message.LogStatusNotificationRequestJson{
		RequestId: &id,
		Status:    message.UploadLogStatusEnumType_1_Idle,
	}
	msg, _, err := message.New("2", "LogStatusNotification", request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) LogStatusNotificationResponse(msgID string, msg []byte) error {
	response := &message.LogStatusNotificationResponseJson{}
	return json.Unmarshal(msg, response)
}
