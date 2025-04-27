package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
	"errors"
)

type Channel struct {
	Main            chan event.Event
	ProcessedEvents chan event.Event
	Funds           chan event.Event
	Orders          chan event.Event
	Trades          chan event.Event
}

type Handler struct {
	MaterializedView *view.MaterializedView
	EventStore       *store.Store
	Channel          Channel
	eventsCounter    int64
	funds            fundsHandler
	orders           ordersHandler
	trades           tradesHandler
}

var (
	ErrUnknownEvent     = errors.New("trying to handle unknown event")
	ErrHandlerCaseLogic = errors.New("error after handling the event type")
)

func New(mainChannel chan event.Event, processedEventsChannel chan event.Event, view *view.MaterializedView) Handler {
	var handler Handler
	handler.MaterializedView = view

	handler.Channel.Main = mainChannel
	handler.Channel.ProcessedEvents = processedEventsChannel
	handler.Channel.Funds = make(chan event.Event)
	handler.Channel.Trades = make(chan event.Event)
	handler.Channel.Orders = make(chan event.Event)
	handler.eventsCounter = 0

	fundsHandler := fundsHandler{
		MainChannel:            mainChannel,
		ProcessedEventsChannel: handler.Channel.ProcessedEvents,
		FundsChannel:           handler.Channel.Funds,
		ActiveTrades:           make(map[int64]event.Event),
		MaterializedView:       handler.MaterializedView,
	}

	tradesHandler := tradesHandler{
		MainChannel:            mainChannel,
		ProcessedEventsChannel: handler.Channel.ProcessedEvents,
		TradesChannel:          handler.Channel.Trades,
		MaterializedView:       handler.MaterializedView,
	}

	ordersHandler := ordersHandler{
		MainChannel:            handler.Channel.Main,
		ProcessedEventsChannel: handler.Channel.ProcessedEvents,
		OrdersChannel:          handler.Channel.Orders,
	}

	handler.funds = fundsHandler
	handler.trades = tradesHandler
	handler.orders = ordersHandler

	return handler
}

func (handler *Handler) Run() error {
	errorChannel := make(chan error)

	// start Funds Handler
	go func() {
		err := handler.funds.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start Orders Handler
	go func() {
		err := handler.orders.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start Trades Handler
	go func() {
		err := handler.trades.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	// start listening to Main channel for unprocessed events
	go func() {
		for currEvent := range handler.Channel.Main {
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
		orderCanceled := currEvent.(*event.OrderCanceled)
		orderCanceled.EventID = handler.eventsCounter
		handler.Channel.Orders <- orderCanceled

	case event.OrdersPlacedEvent:
		orderPlaced := currEvent.(*event.OrderPlaced)
		orderPlaced.EventID = handler.eventsCounter
		handler.Channel.Orders <- orderPlaced

	case event.FundsCreditedEvent:
		fundsCredited := currEvent.(*event.FundsCredited)
		fundsCredited.EventID = handler.eventsCounter
		handler.Channel.Funds <- fundsCredited

	case event.FundsDebitedEvent:
		fundsDebited := currEvent.(*event.FundsDebited)
		fundsDebited.EventID = handler.eventsCounter
		handler.Channel.Funds <- fundsDebited

	case event.TradeExecutedEvent:
		tradeExecuted := currEvent.(*event.TradeExecuted)
		tradeExecuted.EventID = handler.eventsCounter
		handler.Channel.Trades <- tradeExecuted

	default:
		return ErrUnknownEvent

	}

	handler.eventsCounter += 1

	return ErrHandlerCaseLogic
}
