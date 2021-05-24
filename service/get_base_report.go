package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) GetBaseReportResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.GetBaseReportRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	msg, _, err = message.New("3", "GetBaseReport", &message.GetBaseReportResponseJson{Status: message.GenericDeviceModelStatusEnumType_1_Accepted}, msgID)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
