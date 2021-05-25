package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
	"time"
)

func TestChargeStation_UpdateFirmwareResponse(t *testing.T) {
	dateTime := time.Now().Local().Format(time.RFC3339)
	request := &message.UpdateFirmwareRequestJson{
		Firmware: message.FirmwareType{
			InstallDateTime:  &dateTime,
			RetrieveDateTime: dateTime,
		},
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	_, msgID, err := message.New("2", "UpdateFirmware", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "UpdateFirmware", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
