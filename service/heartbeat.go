package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) HeartbeatRequest() ([]byte, error) {
	request := &message.HeartbeatRequestJson{}
	msg, _, err := message.New("2", "Heartbeat", request)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) HeartbeatResponse(msgID string, msg []byte) error {
	response := &message.HeartbeatResponseJson{}
	err := json.Unmarshal(msg, response)
	if err != nil {
		return err
	}
	return nil
}
