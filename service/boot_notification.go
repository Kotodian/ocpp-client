package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

// BootNotificationRequest 该桩第一次初始化的时候发送
func (c *ChargeStation) BootNotificationRequest() ([]byte, error) {
	// 创建request
	request := &message.BootNotificationRequestJson{
		Reason: message.BootReasonEnumType_1_PowerUp,
		CustomData: &message.CustomDataType_2{
			VendorId:      "0000",
			ConnectorType: []string{"2"},
		},
		ChargingStation: message.ChargingStationType{
			SerialNumber: &c.Sn,
			VendorName:   c.VendorName,
			Model:        c.Model,
		},
	}
	// 封装成msg
	msg, _, err := message.New("2", "BootNotification", request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) BootNotificationResponse(msgID string, msg []byte) error {
	response := &message.BootNotificationResponseJson{}
	// 接收ac-ocpp发来的response
	err := json.Unmarshal(msg, response)
	// 解析错误
	if err != nil {
		c.entry.Errorln(err)
		return err
	}
	// 如果状态是rejected就直接返回
	if response.Status == message.RegistrationStatusEnumType_1_Rejected {
		return nil
	}
	// 如果状态是pending, 过段时间重新发送BootNotificationRequest
	if response.Status == message.RegistrationStatusEnumType_1_Pending {
		time.AfterFunc(time.Duration(response.Interval)*time.Second, func() {
			request, err := c.BootNotificationRequest()
			if err != nil {
				c.entry.Errorln(err)
				return
			}
			c.Resend <- request
		})
	} else {
		// 确认Heartbeat的时间间隔
		c.interval = time.Duration(response.Interval) * time.Second
		go func() {
			ticker := time.NewTicker(time.Duration(response.Interval) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-c.stop:
					return
				case <-ticker.C:
					request, err := c.HeartbeatRequest()
					if err != nil {
						c.entry.Errorln(err)
						return
					}
					c.Resend <- request
				}
			}
		}()
	}
	return nil

}
