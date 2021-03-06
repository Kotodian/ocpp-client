package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

// HeartbeatRequest 发送心跳
func (c *ChargeStation) HeartbeatRequest() ([]byte, error) {
	// request
	request := &message.HeartbeatRequestJson{
		StationTime: time.Now().Unix(),
	}
	// 封装成msg
	msg, _, err := message.New("2", "Heartbeat", request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	return msg, nil
}

func (c *ChargeStation) HeartbeatResponse(msgID string, msg []byte) error {
	// 接收ac-ocpp发过来的HeartbeatResponse
	response := &message.HeartbeatResponseJson{}
	err := json.Unmarshal(msg, response)
	if err != nil {
		return err
	}
	return nil
}
