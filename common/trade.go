package common

type Trade struct {
	ID         int64         `json:"id"`
	AccountIDs []int64       `json:"account_id"`
	Orders     []OrderPlaced `json:"orders"`
}

type TradeExecuted struct {
	Trade Trade
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
