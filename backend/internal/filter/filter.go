package filter

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
)

type Filter struct {
	EventStore      store.Store
	View            view.MaterializedView
	processedEvents chan event.Event
	canceledOrders  map[int64]struct{}
}

func New(eventStore store.Store, view view.MaterializedView, processedEvents chan event.Event) *Filter {
	return &Filter{
		EventStore:      eventStore,
		View:            view,
		processedEvents: processedEvents,
		canceledOrders:  make(map[int64]struct{}),
	}
}

func (filter *Filter) Run() error {

	for currEvent := range filter.processedEvents {
		eventType := currEvent.Type()

		var filteredEvent event.Event

		switch eventType {

		case event.OrdersCanceledEvent:
			orderCanceled := currEvent.(*event.OrderCanceled)
			filter.canceledOrders[orderCanceled.Order.ID] = struct{}{}

		case event.OrdersPlacedEvent:
			orderPlaced := currEvent.(*event.OrderPlaced)
			_, exists := filter.canceledOrders[orderPlaced.Order.ID]
			if exists {
				// skip this event
				continue
			}

		case event.FundsCreditedEvent:
			fundsCredited := currEvent.(*event.FundsCredited)
			orderID1 := fundsCredited.Trade.Orders[0].ID
			orderID2 := fundsCredited.Trade.Orders[1].ID
			_, exists1 := filter.canceledOrders[orderID1]
			_, exists2 := filter.canceledOrders[orderID2]
			if exists1 || exists2 {
				// skip this event
				continue
			}

		case event.FundsDebitedEvent:
			fundsDebited := currEvent.(*event.FundsDebited)
			orderID1 := fundsDebited.Trade.Orders[0].ID
			orderID2 := fundsDebited.Trade.Orders[1].ID
			_, exists1 := filter.canceledOrders[orderID1]
			_, exists2 := filter.canceledOrders[orderID2]
			if exists1 || exists2 {
				// skip this event
				continue
			}

		case event.TradeExecutedEvent:
			tradeExecuted := currEvent.(*event.TradeExecuted)
			orderID1 := tradeExecuted.Trade.Orders[0].ID
			orderID2 := tradeExecuted.Trade.Orders[1].ID
			_, exists1 := filter.canceledOrders[orderID1]
			_, exists2 := filter.canceledOrders[orderID2]
			if exists1 || exists2 {
				// skip this event
				continue
			}

		}

		filteredEvent = currEvent
		err := filter.EventStore.AppendEvent(filteredEvent)
		if err != nil {
			return err
		}

		err = filter.View.RegisterEvent(filteredEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
