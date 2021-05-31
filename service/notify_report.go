package service

import (
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) notifyReportRequest(done chan struct{}) {
	msg, _, _ := message.New("2", "NotifyReport", &message.NotifyReportRequestJson{
		GeneratedAt: time.Now().Format(time.RFC3339),
		ReportData: []message.ReportDataType{
			message.ReportDataType{
				Component: message.ComponentType_7{
					Name: "Controller",
				},
				Variable: message.VariableType_6{
					Name: "Manufacturer",
				},
				VariableAttribute:       nil,
				VariableCharacteristics: nil,
			},
		},
		RequestId: 0,
		SeqNo:     0,
		Tbc:       true,
	})
	c.Resend <- msg
	done <- struct{}{}
}
