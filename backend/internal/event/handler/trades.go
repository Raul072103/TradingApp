package handler

import "TradingSimulation/backend/internal/event"

type tradesHandler struct {
}

func (handler *tradesHandler) TradeExecuted(trade event.Trade) error {
	// extra trading logic here
	return nil
}
