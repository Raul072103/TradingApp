package event

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

type EventTypeJSON struct {
	Type int64 `json:"type"`
}

type Event interface {
	Type() int64
	ID() int64
	AccountIDs() []int64
}

// UnmarshalEventTypeJSON unmarshal an EventJSON struct into the event code.
func UnmarshalEventTypeJSON(data []byte) (int64, error) {
	var event EventTypeJSON
	err := json.Unmarshal(data, &event)
	return event.Type, err
}
