package service

import (
	"encoding/json"
	"ocpp-client/message"
)

// FirmwareStatusNotificationRequest 固件状态上报
func (c *ChargeStation) FirmwareStatusNotificationRequest() ([]byte, error) {
	// 随便写的
	requestId := 1
	//请求request
	request := message.FirmwareStatusNotificationRequestJson{
		RequestId: &requestId,
		Status:    message.FirmwareStatusEnumType_1_Idle,
	}
	// 封装成msg
	msg, _, err := message.New("2", "FirmwareStatusNotification", request)
	return msg, err
}

func (c *ChargeStation) FirmwareStatusNotificationResponse(msgID string, msg []byte) error {
	response := &message.FirmwareStatusNotificationResponseJson{}
	// 解析从ac-ocpp发过来的请求
	return json.Unmarshal(msg, response)
}
