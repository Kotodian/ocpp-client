package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) DataTransferResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.DataTransferRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	msg, _, err = message.New("3", "DataTransfer", &message.DataTransferResponseJson{
		Status: message.DataTransferStatusEnumType_1_Accepted,
		Data: &struct {
			CanID     string `json:"canId,omitempty"`
			CanLength string `json:"canLength,omitempty"`
			CanData   string `json:"canData,omitempty"`
		}{
			CanID:     "1231",
			CanLength: "1",
			CanData:   "01",
		},
	}, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}
