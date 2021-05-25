package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_ResetResponse(t *testing.T) {
	request := &message.ResetRequestJson{
		Type: message.ResetEnumType_1_OnIdle,
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, msgID, err := message.New("2", "Reset", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "Reset", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
