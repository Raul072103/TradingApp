package event

import "time"

const (
	SuccessfulTrade = "SUCCESSFUL"
	ActiveTrade     = "ACTIVE"
)

type Trade struct {
	ID         int64   `json:"id"`
	AccountIDs []int64 `json:"account_id"`
	Orders     []Order `json:"orders"`
	Status     string  `json:"status"`
}

type TradeExecuted struct {
	EventID   int64     `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
	Trade     Trade     `json:"trade"`
}

func (e *TradeExecuted) Type() int64 {
	return TradeExecutedEvent
}

func (e *TradeExecuted) ID() int64 {
	return e.Trade.ID
}

func (e *TradeExecuted) AccountIDs() []int64 {
	return e.Trade.AccountIDs
}

func (e *TradeExecuted) Time() time.Time {
	return e.Timestamp
}
