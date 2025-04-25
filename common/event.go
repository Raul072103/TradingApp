package common

const (
	OrdersCanceledEvent = 0
	OrdersPlacedEvent   = 1
	FundsCreditedEvent  = 2
	FundsDebitedEvent   = 3
	TradeExecutedEvent  = 4
)

type Event interface {
	Type() int
	ID() int64
	AccountIDs() []int64
}
