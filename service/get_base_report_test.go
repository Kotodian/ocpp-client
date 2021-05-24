package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_GetBaseReportResponse(t *testing.T) {
	request := &message.GetBaseReportRequestJson{ReportBase: message.ReportBaseEnumType_1_SummaryInventory}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	_, msgID, err := message.New("2", "GetBaseReport", request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "GetBaseReport", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
