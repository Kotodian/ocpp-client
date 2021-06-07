package initialize

import (
	cmap "github.com/orcaman/concurrent-map"
	"ocpp-client/websocket"
)

// 初始化缓存
func init() {
	websocket.Cache = cmap.New()
}
