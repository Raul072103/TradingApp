package handler

import "TradingSimulation/backend/internal/event"

type ordersHandler struct {
	MainChannel            chan event.Event
	ProcessedEventsChannel chan event.Event
	OrdersChannel          chan event.Event
}

func (handler *ordersHandler) Run() error {
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

	return nil
}

func (handler *ordersHandler) handleCanceledOrder(currEvent event.Event) (event.Event, error) {
	return nil, nil
}

func (handler *ordersHandler) handlePlacedOrder(currEvent event.Event) (event.Event, error) {
	return nil, nil
}
