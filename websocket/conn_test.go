package websocket

import (
	"github.com/stretchr/testify/assert"
	"net/url"
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
	client := NewClient(service.NewChargeStation("T1641735200"))
	host := "10.43.0.72:8844"
	addr := url.URL{Scheme: "ws", Host: host, Path: "/ocpp/T1641735200"}
	err := client.Conn(addr.String())
	assert.Nil(t, err)
	time.Sleep(10 * time.Second)
}
