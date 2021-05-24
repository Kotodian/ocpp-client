package websocket

import (
	"github.com/gorilla/websocket"
	"ocpp-client/message"
	"ocpp-client/service"
	"path"
	"strings"
	"sync"
	"time"
)

type Client struct {
	// 锁
	lock sync.Mutex
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

func NewClient(instance *service.ChargeStation) *Client {
	return &Client{
		dialer:   websocket.DefaultDialer,
		instance: instance,
		write:    make(chan []byte, 100),
		read:     make(chan []byte, 100),
	}
}

func (c *Client) Conn(addr string) error {
	// 路由参数
	addr = path.Join(addr, c.instance.ID())
	conn, _, err := c.dialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	if c.instance == nil {
		return nil
	}
	c.conn = conn
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.instance.Stop()
		time.AfterFunc(180*time.Second, func() {
			_ = c.reConn(addr)
		})
		return nil
	})
	go c.writePump()
	go c.readPump()
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
		time.AfterFunc(180*time.Second, func() {
			_ = c.reConn(addr)
		})
		return nil
	})
	c.close = make(chan struct{}, 1)
	go c.writePump()
	go c.readPump()
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
			switch typ {
			case websocket.PingMessage:
				err = c.conn.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					return
				}
			case websocket.TextMessage:
				for _, m := range strings.Split(string(msg), "\n") {
					messageType, messageID, action, payload := message.Parse([]byte(m))
					msg, err = c.instance.Function(messageType, messageID, action, payload)
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
