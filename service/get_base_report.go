package service

import (
	"encoding/json"
	"ocpp-client/message"
)

// GetBaseReportResponse 获取基本报告
func (c *ChargeStation) GetBaseReportResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.GetBaseReportRequestJson{}
	// 解析从ac-ocpp发过来的请求
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	// 封装msg
	msg, _, err = message.New("3", "GetBaseReport", &message.GetBaseReportResponseJson{Status: message.GenericDeviceModelStatusEnumType_1_Accepted}, msgID)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	// 异步地发送多个NotifyReport
	//go func() {
	//	done := make(chan struct{})
	//	c.notifyReportRequest(done)
	//	<-done
	//}()
	//return msg, nil
	return nil, nil
}
