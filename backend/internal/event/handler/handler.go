package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
	"errors"
)

type Channel struct {
	MainChannel chan event.Event
	Funds       chan event.Event
	Orders      chan event.Event
	Trades      chan event.Event
}

type Handler struct {
	MaterializedView *view.MaterializedView
	EventStore       *store.Store
	Channel          Channel
}

var (
	ErrUnknownEvent      = errors.New("trying to handle unknown event")
	ErrInsufficientFunds = errors.New("the user doesn't have sufficient funds")
	ErrHandlerCaseLogic  = errors.New("error after handling the event type")
)

func Run(mainChannel chan event.Event, view *view.MaterializedView) error {
	var handler Handler
	handler.MaterializedView = view

	handler.Channel.MainChannel = mainChannel
	handler.Channel.Funds = make(chan event.Event)
	handler.Channel.Trades = make(chan event.Event)
	handler.Channel.Orders = make(chan event.Event)

	errorChannel := make(chan error)

	fundsHandler := fundsHandler{
		MainChannel:      mainChannel,
		FundsChannel:     handler.Channel.Funds,
		ActiveTrades:     make(map[int64]event.Trade),
		MaterializedView: handler.MaterializedView,
	}

	tradesHandler := tradesHandler{
		MainChannel:      mainChannel,
		TradesChannel:    handler.Channel.Trades,
		MaterializedView: handler.MaterializedView,
	}

	ordersHandler := ordersHandler{
		OrdersChannel: handler.Channel.Orders,
	}

	// start Funds Handler
	go func() {
		err := fundsHandler.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start Orders Handler
	go func() {
		err := ordersHandler.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start Trades Handler
	go func() {
		err := tradesHandler.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start listening to Main channel
	go func() {
		for currEvent := range handler.Channel.MainChannel {
			err := handler.HandleEvent(currEvent)
			if err != nil {
				errorChannel <- err
			}
		}
	}()

	return <-errorChannel
}

func (handler *Handler) HandleEvent(currEvent event.Event) error {
	eventType := currEvent.Type()

	switch eventType {

	case event.OrdersCanceledEvent:
		handler.Channel.Orders <- currEvent

	case event.OrdersPlacedEvent:
		handler.Channel.Orders <- currEvent

	case event.FundsCreditedEvent:
		handler.Channel.Funds <- currEvent

	case event.FundsDebitedEvent:
		handler.Channel.Funds <- currEvent

	case event.TradeExecutedEvent:
		handler.Channel.Trades <- currEvent

	default:
		return ErrUnknownEvent

	}

	return ErrHandlerCaseLogic
}
