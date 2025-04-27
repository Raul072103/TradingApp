package event

import "time"

type FundsDebited struct {
	EventID   int64     `json:"event_id"`
	AccountID int64     `json:"account_id"`
	Timestamp time.Time `json:"timestamp"`
	Sum       float64   `json:"sum"`
	Trade     Trade     `json:"trade"`
}

type FundsCredited struct {
	EventID   int64     `json:"event_id"`
	AccountID int64     `json:"account_id"`
	Timestamp time.Time `json:"timestamp"`
	Sum       float64   `json:"sum"`
	Trade     Trade     `json:"trade"`
}

func (e *FundsDebited) Type() int64 {
	return FundsDebitedEvent
}

func (e *FundsDebited) ID() int64 {
	return e.EventID
}

func (e *FundsDebited) AccountIDs() []int64 {
	return []int64{e.AccountID}
}

func (e *FundsDebited) Time() time.Time {
	return e.Timestamp
}

func (e *FundsCredited) Type() int64 {
	return FundsCreditedEvent
}

func (e *FundsCredited) ID() int64 {
	return e.EventID
}

func (e *FundsCredited) AccountIDs() []int64 {
	return []int64{e.AccountID}
}

func (e *FundsCredited) Time() time.Time {
	return e.Timestamp
}
