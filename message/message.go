package message

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	cache sync.Map
)

// New 新建一个发送给ac-ocpp的message
// 格式如下 [2,"messageID", "action",payload] 或者 [3, "messageID", payload]
// request, _ := json.Marshal(&BootNotificationRequest{})
// New("2", "BootNotification",request)
// response, _ := json.Marshal(&BootNotificationResponse{})
// New("3", "Reset", response, messageID)
func New(typ, action string, payload interface{}, id ...string) ([]byte, string, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, "", nil
	}
	var msgID string
	if typ == "2" {
		msgID = messageID()
		cache.Store(msgID, action)
	} else {
		_, ok := cache.Load(id[0])
		if !ok {
			return nil, "", nil
		}
		msgID = id[0]
	}

	var builder strings.Builder
	builder.WriteByte('[')
	// message type
	builder.WriteString(typ)
	builder.WriteByte(',')
	// message id
	builder.WriteByte('"')
	builder.WriteString(msgID)
	builder.WriteByte('"')
	builder.WriteByte(',')
	if typ == "2" {
		// action
		builder.WriteByte('"')
		builder.WriteString(action)
		builder.WriteByte('"')
		builder.WriteByte(',')
	}
	// payload
	builder.Write(p)

	builder.WriteByte(']')
	return []byte(builder.String()), msgID, nil
}

// messageID 目前是使用时间戳
func messageID() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// Parse 解析ac-ocpp过来的message
func Parse(msg []byte) (typ string, msgID string, action string, payload []byte) {
	// 如果是请求的话
	typ = gjson.GetBytes(msg, "0").String()
	if typ == "2" {
		results := gjson.GetManyBytes(msg, "1", "2", "3")
		msgID = results[0].String()
		action = results[1].String()
		payload = []byte(results[2].Raw)
		cache.Store(msgID, action)
	} else if typ == "3" {
		results := gjson.GetManyBytes(msg, "1", "2")
		msgID = results[0].String()
		a, ok := cache.Load(msgID)
		if !ok {
			return
		}
		payload = []byte(results[1].Raw)

		action = a.(string)
	}
	return
}
