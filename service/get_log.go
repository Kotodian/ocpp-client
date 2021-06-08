package service

import (
	"encoding/json"
	"ocpp-client/message"
)

// GetLogResponse 获取日志
func (c *ChargeStation) GetLogResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.GetLogRequestJson{}
	// 解析平台那边发过来的请求
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	// response
	response := &message.GetLogResponseJson{
		Status: message.LogStatusEnumType_1_Accepted,
	}
	// 封装msg
	msg, _, err = message.New("3", "GetLog", response, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}
