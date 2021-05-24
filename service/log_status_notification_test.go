package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_LogStatusNotificationRequest(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "LogStatusNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestChargeStation_LogStatusNotificationResponse(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "LogStatusNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, id, _, _ := message.Parse(msg)
	assert.NotEmpty(t, id)
	msg, err = json.Marshal(&message.LogStatusNotificationResponseJson{})
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, err = chargeStation.Function("3", id, "LogStatusNotification", msg)
	assert.Nil(t, err)
}
