package common

type Trade struct {
	ID         int64
	AccountIDs []int64
}

type TradeExecuted struct {
	Trade Trade
}

func (e *TradeExecuted) Type() string {
	return TradeExecutedEvent
}

func (e *TradeExecuted) ID() int64 {
	return e.Trade.ID
}

func (e *TradeExecuted) AccountIDs() []int64 {
	return e.Trade.AccountIDs
}
