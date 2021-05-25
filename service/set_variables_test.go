package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_SetVariablesResponse(t *testing.T) {
	request := &message.SetVariablesRequestJson{
		SetVariableData: []message.SetVariableDataType{message.SetVariableDataType{}},
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, msgID, err := message.New("2", "SetVariables", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "SetVariables", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
