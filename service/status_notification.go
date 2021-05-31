package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) StatusNotificationRequest() ([]byte, error) {
	request := &message.StatusNotificationRequestJson{
		ConnectorId:     c.connectors[0].ID(),
		ConnectorStatus: c.connectors[0].State(),
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
