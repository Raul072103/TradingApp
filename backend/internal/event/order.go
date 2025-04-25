package event

import "errors"

const (
	BuyOrder  = "BUY_ORDER"
	SellOrder = "SELL_ORDER"
)

var ErrOrderTypeIncorrect = errors.New("the order type can be either Buy or Sell")

type Order struct {
	ID        int64
	AccountID int64
	Type      string
}

func NewOrder(orderType string) (Order, error) {
	var order Order

	if orderType != BuyOrder && orderType != SellOrder {
		return order, ErrOrderTypeIncorrect
	}

	return order, nil
}

type OrderPlaced struct {
	Order Order
}

type OrderCanceled struct {
	Order Order
}

func (e *OrderPlaced) Type() string {
	return OrdersPlacedEvent
}

func (e *OrderPlaced) ID() int64 {
	return e.Order.ID
}

func (e *OrderPlaced) AccountIDs() []int64 {
	return []int64{e.Order.AccountID}
}

func (e *OrderCanceled) Type() string {
	return OrdersCanceledEvent
}

func (e *OrderCanceled) ID() int64 {
	return e.Order.ID
}

func (e *OrderCanceled) AccountIDs() []int64 {
	return []int64{e.Order.AccountID}
}
