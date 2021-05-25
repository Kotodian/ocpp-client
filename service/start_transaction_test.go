package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"ocpp-client/message"
	"testing"
)

func TestChargeStation_RequestRemoteStartTransactionResponse(t *testing.T) {
	transactionId := "1"
	minChargingRate := 8.1
	period := message.ChargingSchedulePeriodType_4{
		Limit: 8.1,
	}
	schedule := message.ChargingScheduleType_3{
		ChargingRateUnit: message.ChargingRateUnitEnumType_10_A,
		Id:               1,
		MinChargingRate:  &minChargingRate,
		ChargingSchedulePeriod: []message.ChargingSchedulePeriodType_4{
			period,
		},
	}
	// 建立请求
	request := &message.RequestStartTransactionRequestJson{
		ChargingProfile: &message.ChargingProfileType_1{
			ChargingProfileKind:    message.ChargingProfileKindEnumType_3_Absolute,
			ChargingProfilePurpose: message.ChargingProfilePurposeEnumType_7_TxDefaultProfile,
			ChargingSchedule:       []message.ChargingScheduleType_3{schedule},
			TransactionId:          &transactionId,
		},
		RemoteStartId: 1,
		IdToken:       message.IdTokenType_3{Type: message.IdTokenEnumType_7_Local},
	}
	msg, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	// 封装请求
	_, msgID, err := message.New("2", "RequestStartTransaction", request)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgID)
	assert.NotNil(t, msg)

	//处理请求
	chargeStation := NewChargeStation("T164173520")
	msg, err = chargeStation.Function("3", msgID, "RequestStartTransaction", msg)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}
