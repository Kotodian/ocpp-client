package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
	"time"
)

func TestChargeStation_HeartbeatRequest(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "Heartbeat")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestChargeStation_HeartbeatResponse(t *testing.T) {
	chargeStation := NewChargeStation("T164173520")
	msg, err := chargeStation.Function("2", "", "Heartbeat")
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, id, _, _ := message.Parse(msg)
	msg, err = json.Marshal(&message.HeartbeatResponseJson{CurrentTime: time.Now().Format(time.RFC3339)})
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, err = chargeStation.Function("3", id, "Heartbeat", msg)
	assert.Nil(t, err)
}
