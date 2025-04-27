package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/view"
)

type tradesHandler struct {
	MainChannel      chan event.Event
	TradesChannel    chan event.Event
	MaterializedView *view.MaterializedView
}

func (handler *tradesHandler) Run() error {

	return nil
}

func (handler *tradesHandler) TradeExecuted(trade event.Trade) error {
	// extra trading logic here
	return nil
}
