package handler

import "TradingSimulation/backend/internal/event"

type ordersHandler struct {
	MainChannel   chan event.Event
	OrdersChannel chan event.Event
}

func (handler *ordersHandler) Run() error {

	return nil
}

func (handler *ordersHandler) CancelOrder(order event.Order) error {
	// extra canceling order here
	return nil
}

func (handler *ordersHandler) PlacedOrder(order event.Order) error {
	// TODO() check if there is a match
	return nil
}
