package event

const (
	OrdersCanceledEvent = "ORDERS_CANCELED_EVENT"
	OrdersPlacedEvent   = "ORDERS_PLACED_EVENT"
	FundsCreditedEvent  = "FUNDS_CREDITED_EVENT"
	FundsDebitedEvent   = "FUNDS_DEBITED_EVENT"
	TradeExecutedEvent  = "TRADE_EXECUTED_EVENT"
)

type Event interface {
	GetType() string
	GetID() int64
	GetAccountID() int64
}
