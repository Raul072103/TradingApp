package event

import "time"

const (
	BuyOrder  = "BUY"
	SellOrder = "SELL"
)

type Order struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"` // BuyOrder || SellOrder
	Count int64  `json:"count"`
	Stock int64  `json:"stock"`
}

type OrderPlaced struct {
	EventID   int64     `json:"event_id"`
	AccountID int64     `json:"account_id"`
	Timestamp time.Time `json:"timestamp"`
	Order     Order     `json:"order"`
}

type OrderCanceled struct {
	EventID   int64     `json:"event_id"`
	AccountID int64     `json:"account_id"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason"`
	Order     Order     `json:"order"`
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

func (e *OrderPlaced) Time() time.Time {
	return e.Timestamp
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

func (e *OrderCanceled) Time() time.Time {
	return e.Timestamp
}
