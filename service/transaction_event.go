package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"ocpp-client/message"
	"strconv"
	"strings"
	"time"
)

const (
	minElectricity = 1.00
	maxElectricity = 60.00
)

// TransactionEventRequest 如果状态是Started的话就发一次,
// 如果是Updated且正在充电就一直发直到停止充电为止或者连接断开为止
func (c *ChargeStation) TransactionEventRequest() ([]byte, error) {
	if c.Transaction == nil {
		return nil, nil
	}
	if c.Transaction.EventType == message.TransactionEventEnumType_1_Started ||
		c.Transaction.EventType == message.TransactionEventEnumType_1_Ended {
		request := &message.TransactionEventRequestJson{
			EventType: c.Transaction.EventType,
			IdToken: &message.IdTokenType_6{
				IdToken: c.Transaction.IdToken.IdToken,
				Type:    message.IdTokenEnumType_13(c.Transaction.IdTokenType),
			},
			Timestamp:       time.Now().Format(time.RFC3339),
			TransactionInfo: *c.Transaction.Instance,
			SeqNo:           c.Transaction.SeqNo,
		}
		meterValue := genMeterValue(c.Transaction.EventType)
		request.MeterValue = meterValue
		msg, _, err := message.New("2", "TransactionEvent", request)
		return msg, err
	} else if c.Transaction.EventType == message.TransactionEventEnumType_1_Updated {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-c.stop:
					return
				case <-c.Transaction.stop:
					return
				case <-ticker.C:
					c.lock.Lock()
					if c.Transaction.EventType == message.TransactionEventEnumType_1_Ended {
						c.lock.Unlock()
						return
					}
					// 发送meter value
					c.Transaction.Next()
					request := &message.TransactionEventRequestJson{
						EventType: message.TransactionEventEnumType_1_Updated,
						IdToken: &message.IdTokenType_6{
							IdToken: c.Transaction.IdToken.IdToken,
							Type:    message.IdTokenEnumType_13(c.Transaction.IdTokenType),
						},
						SeqNo:           c.Transaction.SeqNo,
						Timestamp:       time.Now().Format(time.RFC3339),
						TransactionInfo: *c.Transaction.Instance,
					}
					// 自动增加电量
					c.Electricity += genElectricity()
					// 充满了
					if c.Electricity >= maxElectricity {
						c.Transaction.entry.Infof("%s 充电完成 电量: %f\n", c.Sn, maxElectricity)
						c.Transaction.EventType = message.TransactionEventEnumType_1_Ended
						c.Electricity = minElectricity
						request.EventType = message.TransactionEventEnumType_1_Ended
						request.TriggerReason = message.TriggerReasonEnumType_1_EnergyLimitReached
						request.MeterValue = genMeterValue(c.Transaction.EventType)
						msg, _, _ := message.New("2", "TransactionEvent", request)
						c.Resend <- msg
						time.Sleep(1 * time.Second)
						c.Connectors[0].SetState(message.ConnectorStatusEnumType_1_Available)
						msg, _ = c.StatusNotificationRequest()
						c.Resend <- msg
						_ = DB.Put(c.BucketName(), c.ID(), c)
						c.lock.Unlock()
						return
					}
					c.Transaction.entry.Infof("%s 充电中 电量: %f\n", c.Sn, c.Electricity)
					request.MeterValue = genMeterValue(c.Transaction.EventType, c.Electricity)
					msg, _, _ := message.New("2", "TransactionEvent", request)
					c.Resend <- msg
					_ = DB.Put(c.BucketName(), c.ID(), c)
					c.lock.Unlock()
				}
			}
		}()
	}
	return nil, nil
}

func (c *ChargeStation) TransactionEventResponse(msgID string, msg []byte) error {
	response := &message.TransactionEventResponseJson{}
	return json.Unmarshal(msg, response)
}

func (c *ChargeStation) SendEvent() {
	c.Transaction.Next()
	msg, _ := c.TransactionEventRequest()
	c.Resend <- msg
}

// genMeterValue 生成MeterValue
func genMeterValue(eventType message.TransactionEventEnumType_1, electricity ...float64) []message.MeterValueType_1 {
	// 包含采样值、上下文以及单位以及数据类型 如果是开始的话提供个初始电量即可(即0.0)
	meterValues := make([]message.MeterValueType_1, 0)
	// 采样值
	sampleValues := make([]message.SampledValueType_1, 0)
	if eventType == message.TransactionEventEnumType_1_Started {
		context := message.ReadingContextEnumType_2_TransactionBegin
		measurand := message.MeasurandEnumType_3_EnergyActiveImportRegister
		// 初始电量为0.0
		sampleValue := message.SampledValueType_1{
			Context:   &context,
			Measurand: &measurand,
			UnitOfMeasure: &message.UnitOfMeasureType_1{
				Unit: "Wh",
			},
			Value: minElectricity,
		}
		sampleValues = append(sampleValues, sampleValue)
		// 加入到meterValues
		meterValues = append(meterValues, message.MeterValueType_1{SampledValue: sampleValues})
	} else if eventType == message.TransactionEventEnumType_1_Updated {
		// 提供固定的功率、电压、电流 不断变化的电量 并且meterValue仍旧只有一个
		context := message.ReadingContextEnumType_2_SampleClock
		// 电量 在之前的电量基础上加
		electricityMeasurand := message.MeasurandEnumType_3_EnergyActiveImportRegister
		sampleValue := message.SampledValueType_1{
			Context:       &context,
			Measurand:     &electricityMeasurand,
			UnitOfMeasure: &message.UnitOfMeasureType_1{Unit: "Wh"},
			Value:         electricity[0],
		}
		sampleValues = append(sampleValues, sampleValue)
		// 功率
		powerMeasurand := message.MeasurandEnumType_3_PowerActiveImport
		sampleValue = message.SampledValueType_1{
			Context:       &context,
			Measurand:     &powerMeasurand,
			UnitOfMeasure: &message.UnitOfMeasureType_1{Unit: "W"},
			Value:         100,
		}
		sampleValues = append(sampleValues, sampleValue)
		// 电压
		voltageMeasurand := message.MeasurandEnumType_3_Voltage
		sampleValue = message.SampledValueType_1{
			Context:       &context,
			Measurand:     &voltageMeasurand,
			UnitOfMeasure: &message.UnitOfMeasureType_1{Unit: "V"},
			Value:         100,
		}
		sampleValues = append(sampleValues, sampleValue)
		// 电流
		currentMeasurand := message.MeasurandEnumType_3_CurrentImport
		sampleValue = message.SampledValueType_1{
			Context:       &context,
			Measurand:     &currentMeasurand,
			UnitOfMeasure: &message.UnitOfMeasureType_1{Unit: "A"},
			Value:         100,
		}
		sampleValues = append(sampleValues, sampleValue)
		// 加入到meterValues
		meterValues = append(meterValues, message.MeterValueType_1{SampledValue: sampleValues})
	} else {
		// 如果是ended meterValue有两个
		context := message.ReadingContextEnumType_2_TransactionEnd
		sampleValues = append(sampleValues, message.SampledValueType_1{Context: nil})
		meterValues = append(meterValues, message.MeterValueType_1{SampledValue: sampleValues})
		sampleValues = make([]message.SampledValueType_1, 0)
		sampleValue := message.SampledValueType_1{
			Context: &context,
			UnitOfMeasure: &message.UnitOfMeasureType_1{
				Unit: "Wh",
			},
			Value: maxElectricity,
		}
		sampleValues = append(sampleValues, sampleValue)
		meterValues = append(meterValues, message.MeterValueType_1{SampledValue: sampleValues})
	}
	return meterValues
}

func genElectricity() float64 {
	f := (rand.Float64() * 5) + 5
	return decimal(f)
}

func decimal(value float64) float64 {
	temp := fmt.Sprintf("%.2f", value)
	value, _ = strconv.ParseFloat(temp, 64)
	inte := strings.Split(temp, ".")[0]
	dec := strings.Split(temp, ".")[1]
	if dec[1] == '0' {
		dec = dec[:1] + "1"
	}
	temp = inte + "." + dec
	value, _ = strconv.ParseFloat(temp, 64)
	return value
}
