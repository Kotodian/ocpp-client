package service

import "ocpp-client/message"

type Transaction struct {
	// 具体的参数
	instance *message.TransactionType
	// 事件类型
	eventType message.TransactionEventEnumType_1
}

func NewTransaction(instance *message.TransactionType) *Transaction {
	return &Transaction{
		instance:  instance,
		eventType: message.TransactionEventEnumType_1_Started,
	}
}
