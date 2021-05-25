package service

import (
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_FirmwareStatusNotificationRequestResponse(t *testing.T) {
	chargeStation := NewChargeStation("T164713520")

	msg, err := chargeStation.Function("2", "", "FirmwareStatusNotification")
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, id, _, payload := message.Parse(msg)
	_, err = chargeStation.Function("3", id, "FirmwareStatusNotification", payload)
	assert.Nil(t, err)
}
