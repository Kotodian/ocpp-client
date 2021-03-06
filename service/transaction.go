package service

import (
	"github.com/sirupsen/logrus"
	"ocpp-client/log"
	"ocpp-client/message"
)

type Transaction struct {
	// 日志
	entry *logrus.Entry
	// 具体的参数
	Instance *message.TransactionType `json:"instance"`
	// 事件类型
	EventType message.TransactionEventEnumType_1 `json:"event_type"`
	// 自增序列号
	SeqNo int `json:"seq_no"`
	// IdToken
	IdToken *message.IdTokenType_3 `json:"id_token"`
	// token type
	IdTokenType message.IdTokenEnumType_7 `json:"id_token_type"`
	ReserveId   int                       `json:"reserve_id"`
	// 停止充电channel stop_transaction和transaction_event交互使用
	stop chan struct{}
}

func NewTransaction(instance *message.TransactionType) *Transaction {
	transaction := &Transaction{
		Instance:  instance,
		EventType: message.TransactionEventEnumType_1_Started,
		stop:      make(chan struct{}),
		SeqNo:     0,
		entry:     log.NewEntry(),
	}
	defer transaction.withID(instance.TransactionId)
	return transaction
}

// Next 每次发送Transaction都要自增该字段
func (t *Transaction) Next() {
	t.SeqNo += 1
}
