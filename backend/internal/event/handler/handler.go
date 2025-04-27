package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
	"errors"
)

type Handler struct {
	Funds            fundsHandler
	Orders           ordersHandler
	Trades           tradesHandler
	MaterializedView *view.MaterializedView
	EventStore       *store.Store
}

var (
	ErrUnknownEvent     = errors.New("trying to handle unknown event")
	ErrHandlerCaseLogic = errors.New("error after handling the event type")
)

func New() (Handler, error) {
	var handler Handler

	return handler, nil
}

func (handler *Handler) HandleEvent(currEvent event.Event) (event.Event, error) {
	eventType := currEvent.Type()

	switch eventType {

	case event.OrdersCanceledEvent:
		// pass it to orders handler
		orderCanceled := currEvent.(*event.OrderCanceled)
		err := handler.Orders.CancelOrder(orderCanceled.Order)
		if err != nil {
			return nil, err
		}

	case event.OrdersPlacedEvent:
		// if it is a buy Order check if the user has sufficient funds at the moment
		orderPlaced := currEvent.(*event.OrderPlaced)
		if orderPlaced.Order.Type == event.BuyOrder {
			accountState, err := handler.MaterializedView.GetAccount(orderPlaced.AccountID)
			if err != nil {
				return nil, err
			}

			if accountState.Funds < orderPlaced.Order.
		}
		// pass it to orders handler

	case event.FundsCreditedEvent:
		// if insufficient funds, return a CancerOrderEvent
		// pass it to funds handler

	case event.FundsDebitedEvent:
		// check if user has available funds, if not cancelOrder and return a new event
		// TODO() add motive to Canceled Order Event
		// pass it to funds handler

	case event.TradeExecutedEvent:
		// this happens only if there were enough funds in user's account at the time of placing the order
		// pass it to trades handler

	default:
		return nil, ErrUnknownEvent

	}

	return nil, ErrHandlerCaseLogic
}
