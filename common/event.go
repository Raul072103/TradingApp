package common

import (
	"encoding/json"
)

const (
	OrdersCanceledEvent = 0
	OrdersPlacedEvent   = 1
	FundsCreditedEvent  = 2
	FundsDebitedEvent   = 3
	TradeExecutedEvent  = 4
)

type EventJSON struct {
	Type       int64   `json:"type"`
	ID         int64   `json:"id"`
	AccountIDs []int64 `json:"account_ids"`
}

type Event interface {
	Type() int64
	ID() int64
	AccountIDs() []int64
}

// UnmarshalEventJSON unmarshal an EventJSON struct into a JSON byte slice.
func UnmarshalEventJSON(data []byte) (EventJSON, error) {
	var event EventJSON
	err := json.Unmarshal(data, &event)
	return event, err
}
