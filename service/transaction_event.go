package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) TransactionEventRequest() ([]byte, error) {
	if c.transaction == nil {
		return nil, nil
	}
	request := &message.TransactionEventRequestJson{
		EventType: message.TransactionEventEnumType_1_Ended,
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
		Timestamp: time.Now().Format(time.RFC3339),
		TransactionInfo: message.TransactionType{
			TransactionId: c.transaction.id,
		},
	}
	msg, _, err := message.New("2", "TransactionEvent", request)
	return msg, err
}

func (c *ChargeStation) TransactionEventResponse(msgID string, msg []byte) error {
	response := &message.TransactionEventResponseJson{}
	return json.Unmarshal(msg, response)
}
