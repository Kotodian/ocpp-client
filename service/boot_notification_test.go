package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
	"time"
)

func TestChargeStation_BootNotificationRequest(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "BootNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestChargeStation_BootNotificationResponse(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "BootNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, id, _, _ := message.Parse(msg)
	msg, err = json.Marshal(&message.BootNotificationResponseJson{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    10,
		Status:      message.RegistrationStatusEnumType_1_Rejected,
	})
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, err = chargeStation.Function("3", id, "BootNotification", msg)
	assert.Nil(t, err)
}
