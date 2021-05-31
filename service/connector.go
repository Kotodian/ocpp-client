package service

import "ocpp-client/message"

// Connector 充电枪
type Connector struct {
	// 充电枪序号
	id int
	// 状态
	state message.ConnectorStatusEnumType_1
}

func NewConnector(id int) *Connector {
	return &Connector{
		id:    id,
		state: message.ConnectorStatusEnumType_1_Available,
	}
}

func (c *Connector) ID() int {
	return c.id
}

func (c *Connector) SetState(state message.ConnectorStatusEnumType_1) {
	c.state = state
}

func (c *Connector) State() message.ConnectorStatusEnumType_1 {
	return c.state
}
