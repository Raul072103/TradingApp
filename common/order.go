package common

const (
	BuyOrder  = "BUY"
	SellOrder = "SELL"
)

type Order struct {
	ID        int64
	AccountID int64
	Type      int64
}

type OrderPlaced struct {
	Order Order
}

type OrderCanceled struct {
	Order Order
}

func (e *OrderPlaced) Type() int64 {
	return OrdersPlacedEvent
}

func (e *OrderPlaced) ID() int64 {
	return e.Order.ID
}

func (e *OrderPlaced) AccountIDs() []int64 {
	return []int64{e.Order.AccountID}
}

func (e *OrderCanceled) Type() int64 {
	return OrdersCanceledEvent
}

func (e *OrderCanceled) ID() int64 {
	return e.Order.ID
}

func (e *OrderCanceled) AccountIDs() []int64 {
	return []int64{e.Order.AccountID}
}
