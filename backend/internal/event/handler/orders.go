package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/handler/matching"
	"TradingSimulation/backend/internal/event/view"
)

type ordersHandler struct {
	MainChannel            chan event.Event
	ProcessedEventsChannel chan event.Event
	OrdersChannel          chan event.Event
	Stocks                 map[int64]view.Stock
	matchingChannel        chan event.Order // matchingChannel sends to the matching service a new order to match
	cancelOrdersChannel    chan event.Order // cancelOrdersChannel sends to the matching service a cancel order
}

func (handler *ordersHandler) Run() error {
	handler.matchingChannel = make(chan event.Order, 100)
	handler.cancelOrdersChannel = make(chan event.Order, 100)

	matchingService := matching.New(
		handler.MainChannel,
		handler.matchingChannel,
		handler.cancelOrdersChannel,
		handler.Stocks,
	)

	errorChannel := make(chan error, 100)

	go func() {
		err := matchingService.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

	for currEvent := range handler.OrdersChannel {
		eventType := currEvent.Type()

		switch eventType {
		case event.OrdersCanceledEvent:
			processedEvent, err := handler.handleCanceledOrder(currEvent)
			if err != nil {
				return err
			}
			if processedEvent == nil {
				// do nothing
			} else {
				handler.ProcessedEventsChannel <- processedEvent
			}

		case event.OrdersPlacedEvent:
			processedEvent, err := handler.handlePlacedOrder(currEvent)
			if err != nil {
				return err
			}
			if processedEvent == nil {
				// do nothing
			} else {
				handler.ProcessedEventsChannel <- processedEvent
			}
		}
	}

	return <-errorChannel
}

func (handler *ordersHandler) handleCanceledOrder(currEvent event.Event) (event.Event, error) {
	canceledOrder := currEvent.(*event.OrderCanceled).Order
	handler.cancelOrdersChannel <- canceledOrder
	return currEvent, nil
}

func (handler *ordersHandler) handlePlacedOrder(currEvent event.Event) (event.Event, error) {
	placedOrder := currEvent.(*event.OrderPlaced).Order
	handler.matchingChannel <- placedOrder
	return currEvent, nil
}
