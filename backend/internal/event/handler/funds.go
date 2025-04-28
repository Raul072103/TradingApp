package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/view"
)

type fundsHandler struct {
	MainChannel            chan event.Event
	FundsChannel           chan event.Event
	ProcessedEventsChannel chan event.Event
	// Every active trade has 2 events, debited, credited. This map keeps the first funds event that happened.
	ActiveTrades     map[int64]event.Event
	MaterializedView *view.MaterializedView
}

func (handler *fundsHandler) Run() error {

	for currEvent := range handler.FundsChannel {
		eventType := currEvent.Type()

		switch eventType {
		case event.FundsCreditedEvent:
			fundsCredited := currEvent.(*event.FundsCredited)

			events, err := handler.handleFundsCredited(currEvent)
			if err != nil {
				return err
			}
			if events == nil {
				// do nothing
			} else {
				underlyingTrade := fundsCredited.Trade
				underlyingTrade.Status = event.SuccessfulTrade

				executedTradeEvent := event.TradeExecuted{
					Trade: underlyingTrade,
				}

				handler.MainChannel <- &executedTradeEvent
				handler.ProcessedEventsChannel <- events[0]
				handler.ProcessedEventsChannel <- events[1]
			}

		case event.FundsDebitedEvent:
			fundsDebited := currEvent.(*event.FundsDebited)

			events, err := handler.handleFundsDebited(currEvent)
			if err != nil {
				return err
			}
			if events == nil {
				// do nothing
			} else {
				underlyingTrade := fundsDebited.Trade
				underlyingTrade.Status = event.SuccessfulTrade

				executedTradeEvent := event.TradeExecuted{
					Trade:     underlyingTrade,
					AccountID: []int64{underlyingTrade.AccountIDs[0], underlyingTrade.AccountIDs[1]},
				}

				handler.ProcessedEventsChannel <- &executedTradeEvent
				handler.ProcessedEventsChannel <- events[0]
				handler.ProcessedEventsChannel <- events[1]
			}

		}
	}

	return nil
}

func (handler *fundsHandler) handleFundsCredited(currEvent event.Event) ([]event.Event, error) {
	fundsCredited := currEvent.(*event.FundsCredited)

	fundsDebited, exists := handler.ActiveTrades[fundsCredited.Trade.ID]
	if exists {
		fundsDebited = fundsDebited.(*event.FundsDebited)
		delete(handler.ActiveTrades, fundsCredited.Trade.ID)

		return []event.Event{fundsDebited, fundsCredited}, nil
	} else {
		handler.ActiveTrades[fundsCredited.Trade.ID] = fundsCredited
	}

	// for now store the event, until its brother is processed
	return nil, nil
}

func (handler *fundsHandler) handleFundsDebited(currEvent event.Event) ([]event.Event, error) {
	fundsDebited := currEvent.(*event.FundsDebited)

	fundsCredited, exists := handler.ActiveTrades[fundsDebited.Trade.ID]
	if exists {
		fundsCredited = fundsCredited.(*event.FundsCredited)
		delete(handler.ActiveTrades, fundsDebited.Trade.ID)

		return []event.Event{fundsDebited, fundsCredited}, nil
	} else {
		handler.ActiveTrades[fundsDebited.Trade.ID] = fundsDebited
	}

	// for now store the event, until its brother is processed
	return nil, nil
}
