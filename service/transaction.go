package service

import "ocpp-client/message"

type Transaction struct {
	// 具体的参数
	instance *message.TransactionType
	// 事件类型
	eventType message.TransactionEventEnumType_1
	// 自增序列号
	seqNo int
	// IdToken
	idToken *message.IdTokenType_3
	// token type
	idTokenType message.IdTokenEnumType_7
	// 停止充电channel stop_transaction和transaction_event交互使用
	stop chan struct{}
}

func NewTransaction(instance *message.TransactionType) *Transaction {
	return &Transaction{
		instance:  instance,
		eventType: message.TransactionEventEnumType_1_Started,
		stop:      make(chan struct{}),
		seqNo:     0,
	}
}

// Next 每次发送Transaction都要自增该字段
func (t *Transaction) Next() {
	t.seqNo += 1
}
