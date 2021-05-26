package websocket

import (
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"ocpp-client/message"
	"ocpp-client/service"
	"strings"
	"sync"
	"time"
)

var Cache cmap.ConcurrentMap

type Client struct {
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
		close:     make(chan struct{}, 1),
		connected: false,
	}
	Cache.Set(instance.ID(), client)
	return client
}

func (c *Client) Conn(addr string) error {
	// 路由参数
	conn, _, err := c.dialer.Dial(addr, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	if c.instance == nil {
		return nil
	}
	c.conn = conn
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.instance.Stop()
		c.SetConnected(false)
		time.AfterFunc(180*time.Second, func() {
			_ = c.reConn(addr)
		})
		return nil
	})
	go c.writePump()
	go c.readPump()
	time.Sleep(1 * time.Second)
	c.SetConnected(true)
	msg, _ := c.instance.Function("2", "", "BootNotification")
	c.write <- msg
	return nil
}

func (c *Client) reConn(addr string) error {
	conn, _, err := c.dialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.instance.Stop()
		c.SetConnected(false)
		time.AfterFunc(180*time.Second, func() {
			_ = c.reConn(addr)
		})
		return nil
	})
	c.close = make(chan struct{}, 1)
	go c.writePump()
	go c.readPump()
	// 睡一秒保证两个协程建立好
	time.Sleep(1 * time.Second)
	c.SetConnected(true)
	c.instance.ReConn()
	return nil
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.write:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case msg := <-c.instance.Resend():
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer c.conn.Close()
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
