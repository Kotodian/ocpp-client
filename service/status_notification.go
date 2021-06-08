package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) StatusNotificationRequest() ([]byte, error) {
	request := &message.StatusNotificationRequestJson{
		ConnectorId:     c.Connectors[0].ID(),
		ConnectorStatus: c.Connectors[0].State,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
	msg, _, err := message.New("2", "StatusNotification", request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) StatusNotificationResponse(msgID string, msg []byte) error {
	response := &message.StatusNotificationResponseJson{}
	err := json.Unmarshal(msg, response)
	if err != nil {
		c.entry.Errorln(err)
		return err
	}
	return nil
}
