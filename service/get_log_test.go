package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_GetLogResponse(t *testing.T) {
	request := &message.GetLogRequestJson{
		Log: message.LogParametersType{
			RemoteLocation: "/",
		},
		LogType:   message.LogEnumType_1_DiagnosticsLog,
		RequestId: 1,
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, msgID, err := message.New("2", "GetLog", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "GetLog", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
