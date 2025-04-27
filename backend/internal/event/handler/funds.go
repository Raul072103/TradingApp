package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/view"
)

type fundsHandler struct {
	MainChannel  chan event.Event
	FundsChannel chan event.Event
	// Every active trade has 2 events, debited, credited. This map keeps count of the related events based on Trade ID.
	ActiveTrades     map[int64]event.Trade
	MaterializedView *view.MaterializedView
}

func (handler *fundsHandler) Run() error {

	return nil
}

func (handler *fundsHandler) fundsCredited() error {
	// extra funds credited logic here
	return nil
}

func (handler *fundsHandler) fundsDebited() error {
	// extra funds debited logic here
	return nil
}
