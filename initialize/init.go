package initialize

import (
	cmap "github.com/orcaman/concurrent-map"
	"ocpp-client/websocket"
)

func init() {
	websocket.Cache = cmap.New()
	initapi()
	initLog()
	//initdb()
}
