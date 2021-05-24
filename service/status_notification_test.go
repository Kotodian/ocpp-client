package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_StatusNotificationRequest(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "StatusNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestChargeStation_StatusNotificationResponse(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	// 模拟request
	msg, err := chargeStation.Function("2", "", "StatusNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, id, _, _ := message.Parse(msg)
	assert.NotEmpty(t, id)
	msg, err = json.Marshal(&message.StatusNotificationResponseJson{})
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, err = chargeStation.Function("3", id, "StatusNotification", msg)
	assert.Nil(t, err)
}
