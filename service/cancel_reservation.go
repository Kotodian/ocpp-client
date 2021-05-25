package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) CancelReservationResponse(msgID string, msg []byte) ([]byte, error) {
	request := &message.CancelReservationRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.CancelReservationResponseJson{
		Status: message.CancelReservationStatusEnumType_1_Accepted,
	}
	msg, _, err = message.New("3", "CancelReservation", response, msgID)
	return msg, err
}
