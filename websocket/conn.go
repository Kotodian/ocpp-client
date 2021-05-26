package websocket

import (
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	"ocpp-client/message"
	"ocpp-client/service"
	"strings"
	"sync"
	"time"
)

var Cache cmap.ConcurrentMap

type Client struct {
	addr string
	// 锁
	lock sync.Mutex
	// 是否已经连接
	connected bool
	// 客户端
	dialer *websocket.Dialer
	// 连接
	conn *websocket.Conn
	// 充电桩实例
	instance *service.ChargeStation
	// 发送消息
	write chan []byte
	// 接收消息
	read chan []byte
	// 意外关闭
	close chan struct{}
}

func init() {
	Cache = cmap.New()
}

func NewClient(instance *service.ChargeStation) *Client {
	client := &Client{
		lock:      sync.Mutex{},
		dialer:    websocket.DefaultDialer,
		instance:  instance,
		write:     make(chan []byte, 100),
		read:      make(chan []byte, 100),
		connected: false,
	}
	Cache.Set(instance.ID(), client)
	return client
}

func (c *Client) Conn(addr string) error {
	// 路由参数
	conn, _, err := c.dialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	if c.instance == nil {
		return nil
	}
	c.conn = conn
	c.addr = addr
	c.close = make(chan struct{})
	c.SetConnected(true)
	go c.writePump()
	go c.readPump()
	time.Sleep(1 * time.Second)
	msg, _ := c.instance.Function("2", "", "BootNotification")
	c.write <- msg
	return nil
}

func (c *Client) reConn() error {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	c.close = make(chan struct{})
	c.SetConnected(true)
	go c.writePump()
	go c.readPump()
	// 睡一秒保证两个协程建立好
	time.Sleep(1 * time.Second)
	c.instance.ReConn()
	return nil
}

func (c *Client) writePump() {
	defer c.Close()
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.write:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case msg := <-c.instance.Resend:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer c.Close()
	for {
		select {
		case <-c.close:
			return
		default:
			typ, msg, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			if len(msg) == 0 {
				continue
			}
			switch typ {
			case websocket.PingMessage:
				err = c.conn.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					return
				}
			case websocket.TextMessage:
				for _, m := range strings.Split(string(msg), "\n") {
					typ, messageID, action, payload := message.Parse([]byte(m))
					if typ == "4" {
						continue
					}
					msg, err = c.instance.Function("3", messageID, action, payload)
					if err != nil {
						return
					}
					if msg == nil {
						continue
					}
					c.write <- msg
				}
			}
		}
	}
}

func (c *Client) SetConnected(connected bool) {
	c.lock.Lock()
	c.connected = connected
	c.lock.Unlock()
}

func (c *Client) Connected() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.connected
}

func (c *Client) Instance() *service.ChargeStation {
	return c.instance
}

func (c *Client) Write(msg []byte) {
	c.write <- msg
}

func (c *Client) Close() {
	if c.Connected() == false {
		return
	}
	_ = c.conn.Close()
	c.SetConnected(false)
	c.instance.Stop()
	time.AfterFunc(10*time.Second, func() {
		_ = c.reConn()
	})
	close(c.close)
}
