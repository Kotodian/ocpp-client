package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_SendLocalListResponse(t *testing.T) {
	request := &message.SendLocalListRequestJson{
		LocalAuthorizationList: nil,
		UpdateType:             message.UpdateEnumType_1_Full,
		VersionNumber:          0,
	}
	_, msgID, err := message.New("2", "SendLocalList", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)

	chargeStation := NewChargeStation("T164713520")
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	msg, err = chargeStation.Function("3", msgID, "SendLocalList", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
