package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

// TransactionEventRequest 如果状态是Started的话就发一次,
// 如果是Updated且正在充电就一直发直到停止充电为止或者连接断开为止
func (c *ChargeStation) TransactionEventRequest() ([]byte, error) {
	if c.transaction == nil {
		return nil, nil
	}
	request := &message.TransactionEventRequestJson{
		EventType: c.transaction.eventType,
		MeterValue: []message.MeterValueType_1{
			message.MeterValueType_1{
				SampledValue: []message.SampledValueType_1{
					message.SampledValueType_1{
						Context: nil,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
			message.MeterValueType_1{
				SampledValue: []message.SampledValueType_1{
					message.SampledValueType_1{
						Context: nil,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
		Timestamp:       time.Now().Format(time.RFC3339),
		TransactionInfo: *c.transaction.instance,
	}
	msg, _, err := message.New("2", "TransactionEvent", request)
	return msg, err
}

func (c *ChargeStation) TransactionEventResponse(msgID string, msg []byte) error {
	response := &message.TransactionEventResponseJson{}
	return json.Unmarshal(msg, response)
}
