package event

type TradeExecuted struct {
}

func (e *TradeExecuted) Type() string {
	return TradeExecutedEvent
}

func (e *TradeExecuted) ID() int64 {
	return 0
}

func (e *TradeExecuted) AccountID() int64 {
	return 0
}
