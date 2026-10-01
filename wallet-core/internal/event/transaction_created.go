package event

import "time"

type TransactionCreated struct {
	Name    string
	Payload interface{}
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func (tc *TransactionCreated) GetName() string {
	return tc.Name
}

func (tc *TransactionCreated) GetPayload() interface{} {
	return tc.Payload
}

func (tc *TransactionCreated) SetPayload(payload interface{}) {
	tc.Payload = payload
}

func (tc *TransactionCreated) GetDateTime() time.Time {
	return time.Now()
}
