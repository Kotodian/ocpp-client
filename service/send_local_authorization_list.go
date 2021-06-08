package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) SendLocalListResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.SendLocalListRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.SendLocalListResponseJson{
		Status: message.SendLocalListStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "SendLocalList", response, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}
