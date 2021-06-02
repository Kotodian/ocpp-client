package service

import (
	"ocpp-client/message"
	"sync"
)

// Connector 充电枪
type Connector struct {
	// 锁
	lock sync.Mutex
	// 充电枪序号
	Id int `json:"id"`
	// 状态
	State message.ConnectorStatusEnumType_1 `json:"state"`
}

func NewConnector(id int) *Connector {
	return &Connector{
		Id:    id,
		State: message.ConnectorStatusEnumType_1_Available,
	}
}

func (c *Connector) ID() int {
	return c.Id
}

func (c *Connector) SetState(state message.ConnectorStatusEnumType_1) {
	c.State = state
}
