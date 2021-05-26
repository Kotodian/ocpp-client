package service

// Transaction 充电事务
type Transaction struct {
	// 充电事务的id
	id string
}

func NewTransaction(id string) *Transaction {
	return &Transaction{id: id}
}
