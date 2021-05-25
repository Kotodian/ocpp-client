package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_CancelReservationResponse(t *testing.T) {
	request := &message.CancelReservationRequestJson{
		ReservationId: 1,
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, msgID, err := message.New("2", "CancelReservation", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "CancelReservation", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
