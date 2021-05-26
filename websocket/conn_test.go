package websocket

import (
	"github.com/stretchr/testify/assert"
	"ocpp-client/service"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	_ = NewClient(service.NewChargeStation("T164173520"))
	get, _ := Cache.Get("T164173520")
	assert.NotNil(t, get)
}

func TestClient_Conn(t *testing.T) {
	client := NewClient(service.NewChargeStation("T164173520"))
	err := client.Conn("")
	time.Sleep(10 * time.Second)
	assert.Nil(t, err)
}
