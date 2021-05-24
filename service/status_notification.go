package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) StatusNotificationRequest() ([]byte, error) {
	request := &message.StatusNotificationRequestJson{
		ConnectorId:     1,
		ConnectorStatus: message.ConnectorStatusEnumType_1_Available,
		EvseId:          1,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
	msg, _, err := message.New("2", "StatusNotification", request)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) StatusNotificationResponse(msgID string, msg []byte) error {
	response := &message.StatusNotificationResponseJson{}
	err := json.Unmarshal(msg, response)
	if err != nil {
		return err
	}
	return nil
}
