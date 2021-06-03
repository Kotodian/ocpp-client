package websocket

import (
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"ocpp-client/log"
	"ocpp-client/message"
	"ocpp-client/service"
	"strings"
	"sync"
	"time"
)

// Cache 缓存,存储所有的client
var Cache cmap.ConcurrentMap

// Client 桩的websocket客户端
type Client struct {
	// 日志
	entry *logrus.Entry
	// 连接地址
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
	// 关闭
	close chan struct{}
}

// 初始化缓存
func init() {
	Cache = cmap.New()
	err := service.DB.ForEach(service.ChargeStationBucket, func(k string, v interface{}) error {
		chargeStation := v.(*service.ChargeStation)
		err := NewClient(chargeStation).reConn()
		return err
	})
	if err != nil {
		panic(err)
	}
}

// NewClient 新建客户端,根据charge station 实例初始化
func NewClient(instance *service.ChargeStation) *Client {
	client := &Client{
		lock:      sync.Mutex{},
		dialer:    websocket.DefaultDialer,
		instance:  instance,
		write:     make(chan []byte, 100),
		read:      make(chan []byte, 100),
		connected: false,
		entry:     log.NewEntry(),
	}
	defer client.withSN(instance.ID())
	// 存储的缓存
	if _, ok := Cache.Get(instance.ID()); ok {
		return nil
	}
	Cache.Set(instance.ID(), client)
	return client
}

// Conn 连接到指定地址
func (c *Client) Conn(addr string) error {
	// 建立连接
	conn, _, err := c.dialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			return
		}
		c.withAddr(addr)
	}()
	// 如果没有实例 就返回
	if c.instance == nil {
		return nil
	}
	// 设置各个参数
	c.conn = conn
	c.addr = addr
	c.close = make(chan struct{})
	c.SetConnected(true)
	// 开启读写
	go c.writePump()
	go c.readPump()
	// 睡一秒等待协程建立完毕
	time.Sleep(1 * time.Second)
	// 调用BootNotification
	msg, _ := c.instance.Function("2", "", "BootNotification")
	c.write <- msg
	time.Sleep(5 * time.Second)
	// 再次调用BootNotification
	msg, _ = c.instance.Function("2", "", "BootNotification")
	c.write <- msg
	time.Sleep(1 * time.Second)
	msg, _ = c.instance.Function("2", "", "StatusNotification")
	c.write <- msg
	return nil
}

// reConn 重新建立连接
func (c *Client) reConn() error {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			c.withAddr(c.addr)
		}
	}()
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

// writePump 往ac-ocpp发数据
func (c *Client) writePump() {
	defer c.Close()
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.write:
			c.entry.Infoln(msg)
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case msg := <-c.instance.Resend:
			c.entry.Infoln(msg)
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

// readPump 从ac-ocpp读取数据
func (c *Client) readPump() {
	defer c.Close()
	for {
		select {
		case <-c.close:
			return
		default:
			// 读取数据
			typ, msg, err := c.conn.ReadMessage()
			// 有报错就直接关闭连接
			if err != nil {
				return
			}
			// 如果是空数据就不管
			if len(msg) == 0 {
				continue
			}
			// 判断数据类型
			switch typ {
			// 如果ping,就返回pong数据
			case websocket.PingMessage:
				err = c.conn.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					return
				}
			// 如果是普通数据
			case websocket.TextMessage:
				// 可能发过来多条命令以 \n做拆分
				for _, m := range strings.Split(string(msg), "\n") {
					// 记录从ac-ocpp发过来的数据
					c.entry.Infoln(m)
					// 获取 messageType,messageID, action, payload
					typ, messageID, action, payload := message.Parse([]byte(m))
					// 对callError就不管了
					if typ == "4" {
						continue
					}
					// 调用response方法获取response消息
					msg, err = c.instance.Function("3", messageID, action, payload)
					// 如果有错就跳到下个命令
					if err != nil {
						continue
					}
					// 如果数据为空就代表不需要回复到ac-ocpp
					if msg == nil {
						continue
					}
					// 交给writePump
					c.write <- msg
				}
			}
		}
	}
}

// SetConnected 设置是否连接
func (c *Client) SetConnected(connected bool) {
	c.lock.Lock()
	c.connected = connected
	c.lock.Unlock()
}

// Connected 确认是否连接
func (c *Client) Connected() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.connected
}

// Instance 获取charge station实例
func (c *Client) Instance() *service.ChargeStation {
	return c.instance
}

// Write 交给writePump处理
func (c *Client) Write(msg []byte) {
	c.write <- msg
}

// Close 关闭连接
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
