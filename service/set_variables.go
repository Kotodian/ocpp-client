package service

import (
	"encoding/json"
	"ocpp-client/message"
)

func (c *ChargeStation) SetVariables(msgID string, msg []byte) ([]byte, error) {
	request := &message.SetVariablesRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		return nil, err
	}
	response := &message.SetVariablesResponseJson{
		SetVariableResult: []message.SetVariableResultType{
			message.SetVariableResultType{
				AttributeStatus: message.SetVariableStatusEnumTypeAccepted,
			},
		},
	}
	msg, _, err = message.New("3", "SetVariables", response, msgID)
	return msg, err
}
