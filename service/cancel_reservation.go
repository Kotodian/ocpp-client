package service

import (
	"encoding/json"
	"ocpp-client/message"
)

// CancelReservationResponse 取消预约充电
func (c *ChargeStation) CancelReservationResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.CancelReservationRequestJson{}
	// 解析从平台那边发送来的request
	err := json.Unmarshal(msg, request)
	// 解析失败
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	// 返回response
	response := &message.CancelReservationResponseJson{
		Status: message.CancelReservationStatusEnumType_1_Accepted,
	}
	// 封装成msg
	msg, _, err = message.New("3", "CancelReservation", response, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}
