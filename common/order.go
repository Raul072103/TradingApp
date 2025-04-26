package common

const (
	BuyOrder  = "BUY"
	SellOrder = "SELL"
)

type Order struct {
	Type int64 // BuyOrder || SellOrder
}

type OrderPlaced struct {
	EventID   int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	Order     Order `json:"order"`
}

type OrderCanceled struct {
	EventID   int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	Order     Order `json:"order"`
}

func (e *OrderPlaced) Type() int64 {
	return OrdersPlacedEvent
}

func (e *OrderPlaced) ID() int64 {
	return e.EventID
}

func (e *OrderPlaced) AccountIDs() []int64 {
	return []int64{e.AccountID}
}

func (e *OrderCanceled) Type() int64 {
	return OrdersCanceledEvent
}

func (e *OrderCanceled) ID() int64 {
	return e.EventID
}

func (e *OrderCanceled) AccountIDs() []int64 {
	return []int64{e.AccountID}
}
