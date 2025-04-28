package handler

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/view"
)

type tradesHandler struct {
	MainChannel            chan event.Event
	ProcessedEventsChannel chan event.Event
	TradesChannel          chan event.Event
	MaterializedView       *view.MaterializedView
}

func (handler *tradesHandler) Run() error {

	for currEvent := range handler.TradesChannel {
		tradeExecutedEvent := currEvent.(*event.TradeExecuted)
		trade := tradeExecutedEvent.Trade
		stockID := trade.Orders[0].Stock

		if trade.Status == event.SuccessfulTrade {
			handler.ProcessedEventsChannel <- currEvent
		} else {
			// calculate funds
			stock := handler.MaterializedView.Stocks[stockID]

			order1 := trade.Orders[0]
			order2 := trade.Orders[1]

			var debitedAccountID int64
			var creditedAccountID int64

			if order1.Type == event.SellOrder {
				debitedAccountID = order2.AccountID
				creditedAccountID = order1.AccountID
			} else {
				debitedAccountID = order1.AccountID
				creditedAccountID = order2.AccountID
			}

			totalAmount := stock.Price * float64(trade.Orders[0].Count)
			fundsDebited := event.FundsDebited{
				Sum:       totalAmount,
				AccountID: debitedAccountID,
				Trade:     trade,
			}

			fundsCredited := event.FundsCredited{
				Sum:       totalAmount,
				AccountID: creditedAccountID,
				Trade:     trade,
			}

			handler.MainChannel <- &fundsDebited
			handler.MainChannel <- &fundsCredited
		}
	}

	return nil
}
