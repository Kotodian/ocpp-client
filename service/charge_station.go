package service

import (
	"reflect"
	"time"
)

type ChargeStation struct {
	sn         string
	vendorName string
	model      string
	interval   time.Duration
	stop       chan struct{}
	Resend     chan []byte
}

func NewChargeStation(sn string) *ChargeStation {
	return &ChargeStation{
		sn:         sn,
		stop:       make(chan struct{}, 1),
		vendorName: "joysonquin",
		model:      "test",
		Resend:     make(chan []byte, 1),
	}
}

func (c *ChargeStation) ID() string {
	return c.sn
}

func (c *ChargeStation) Function(typ string, messageID, action string, parameters ...interface{}) ([]byte, error) {
	if typ == "2" {
		action += "Request"
	} else {
		action += "Response"
	}

	var p []reflect.Value
	if len(messageID) > 0 {
		p = append(p, reflect.ValueOf(messageID))
	}
	if len(parameters) > 0 {
		for _, parameter := range parameters {
			p = append(p, reflect.ValueOf(parameter))
		}
	}
	values := reflect.ValueOf(c).MethodByName(action).Call(p)
	if len(values) == 2 {
		if values[1].IsNil() {
			return values[0].Bytes(), nil
		} else {
			return nil, values[1].Interface().(error)
		}

	} else {
		if values[0].IsNil() {
			return nil, nil
		} else {
			return nil, values[0].Interface().(error)
		}
	}
}

func (c *ChargeStation) Stop() {
	c.stop <- struct{}{}
}

func (c *ChargeStation) ReConn() {
	if c.interval == 0 {
		msg, _ := c.Function("2", "", "BootNotification")
		c.Resend <- msg
	} else {
		go func() {
			ticker := time.NewTicker(c.interval)
			defer ticker.Stop()
			for {
				select {
				case <-c.stop:
					return
				case <-ticker.C:
					msg, _ := c.Function("2", "", "Heartbeat")
					c.Resend <- msg
				}
			}
		}()
		msg, _ := c.Function("2", "", "StatusNotification")
		c.Resend <- msg
	}
}
