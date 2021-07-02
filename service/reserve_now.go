package service

import (
	"encoding/json"
	"ocpp-client/message"
	"time"
)

func (c *ChargeStation) ReserveNowResponse(msgID string, msg []byte) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	request := &message.ReserveNowRequestJson{}
	err := json.Unmarshal(msg, request)
	if err != nil {
		c.entry.Errorln(err)
		return nil, err
	}
	response := &message.ReserveNowResponseJson{}
	defer func() {
		if err != nil {
			c.entry.Errorln(err)
			return
		}
	}()
	afterTime := getTime(time.Now())
	if request.CustomData.ChargeTime <= afterTime {
		if c.InTransaction() {
			response.Status = message.ReserveNowStatusEnumType_1_Rejected
			goto send
		}
	}
	defer func() {
		time.AfterFunc(time.Duration(afterTime)+1*time.Second, func() {
			c.lock.Lock()
			defer c.lock.Unlock()
			if c.InTransaction() {
				return
			}
			_ = c.StartTransaction()
			time.Sleep(1 * time.Second)
			// 再发送TransactionEventRequest
			c.Transaction.IdTokenType = message.IdTokenEnumType_7_Central
			c.Transaction.IdToken = &message.IdTokenType_3{
				IdToken: request.IdToken.IdToken,
				Type:    message.IdTokenEnumType_7_Central,
			}
			c.Transaction.ReserveId = int(afterTime)
			c.SendEvent()
			time.Sleep(1 * time.Second)
			// 存进数据库
			//err = DB.Put(c.ID(), c)
			if err != nil {
				return
			}
			err = c.UpdateTransaction()
			if err != nil {
				return
			}
		})
	}()
	response.Status = message.ReserveNowStatusEnumType_1_Accepted
send:
	msg, _, _ = message.New("3", "ReserveNow", response, msgID)
	c.Resend <- msg
	return nil, nil
}

func getTime(chargeTime time.Time) int32 {
	return int32(chargeTime.Hour()*3600 + chargeTime.Minute()*60)
}
