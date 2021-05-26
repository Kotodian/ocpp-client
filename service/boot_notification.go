package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) BootNotificationRequest() ([]byte, error) {
	request := &message.BootNotificationRequestJson{
		Reason: message.BootReasonEnumType_1_PowerUp,
		CustomData: &message.CustomDataType_2{
			VendorId:      "0000",
			ConnectorType: []string{"2"},
		},
		ChargingStation: message.ChargingStationType{
			SerialNumber: &c.sn,
			VendorName:   c.vendorName,
			Model:        c.model,
		},
	}
	msg, _, err := message.New("2", "BootNotification", request)
	return msg, err
}

func (c *ChargeStation) BootNotificationResponse(msgID string, msg []byte) error {
	response := &message.BootNotificationResponseJson{}
	err := json.Unmarshal(msg, response)
	if err != nil {
		return err
	}
	if response.Status == message.RegistrationStatusEnumType_1_Rejected {
		return nil
	}

	if response.Status == message.RegistrationStatusEnumType_1_Pending {
		go func() {
			ticker := time.NewTicker(time.Duration(response.Interval) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-c.stop:
					return
				case <-ticker.C:
					request, err := c.BootNotificationRequest()
					if err != nil {
						return
					}
					c.resend <- request
				}
			}
		}()
	} else {
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
						return
					}
					c.resend <- request
				}
			}
		}()
	}

	return nil

}
