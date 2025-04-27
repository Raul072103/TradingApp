package view

import (
	"TradingSimulation/backend/internal/event"
)

// List of accounts
// List of trades
// List of orders

type AccountState struct {
	BuyOrders      map[int64]event.Order
	SellOrders     map[int64]event.Order
	CanceledOrders map[int64]event.Order

	Trades map[int64]event.Trade

	Events []event.Event
}

type MaterializedView struct {
	Accounts map[int64]AccountState
	Trades   []event.Trade
}

// type OrderCanceled struct {
//	EventID   int64 `json:"id"`
//	AccountID int64 `json:"account_id"`
//	Order     Order `json:"order"`
//}

//type FundsDebited struct {
//	EventID   int64   `json:"id"`
//	AccountID int64   `json:"account_id"`
//	Sum       float64 `json:"sum"`
//	Trade     Trade   `json:"trade"`
//}

//type Trade struct {
//	ID         int64         `json:"id"`
//	AccountIDs []int64       `json:"account_id"`
//	Orders     []OrderPlaced `json:"orders"`
//}
//
//type TradeExecuted struct {
//	Trade Trade
//}
