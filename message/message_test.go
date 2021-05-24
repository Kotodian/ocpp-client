package message

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	msg, msgID, err := New("2", "BootNotification",
		&BootNotificationRequestJson{ChargingStation: ChargingStationType{Model: "1"}, Reason: BootReasonEnumType_1_PowerUp})
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	msg, _, err = New("3", "BootNotification", &BootNotificationResponseJson{CurrentTime: time.Now().Format(time.RFC3339), Interval: 10, Status: RegistrationStatusEnumType_1_Accepted}, msgID)
	assert.Nil(t, err)
	t.Log(string(msg))
}

func TestParse(t *testing.T) {
	msg, msgID, err := New("2", "BootNotification",
		&BootNotificationRequestJson{ChargingStation: ChargingStationType{Model: "1"}, Reason: BootReasonEnumType_1_PowerUp})
	assert.Nil(t, err)
	_, _, action, payload := Parse(msg)
	assert.NotEmpty(t, action)
	assert.NotNil(t, payload)
	msg, _, err = New("3", "BootNotification", &BootNotificationResponseJson{CurrentTime: time.Now().Format(time.RFC3339), Interval: 10, Status: RegistrationStatusEnumType_1_Accepted}, msgID)
	assert.Nil(t, err)
	_, _, action, payload = Parse(msg)
	assert.NotEmpty(t, action)
	assert.NotNil(t, payload)
	err = json.Unmarshal(payload, &BootNotificationResponseJson{})
	assert.Nil(t, err)
}
