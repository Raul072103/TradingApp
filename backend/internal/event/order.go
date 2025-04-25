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
	order Order
}

type OrderCanceled struct {
	order Order
}

func (e *OrderPlaced) Type() string {
	return e.order.Type
}

func (e *OrderPlaced) ID() int64 {
	return e.order.ID
}

func (e *OrderPlaced) AccountID() int64 {
	return e.order.AccountID
}

func (e *OrderCanceled) Type() string {
	return e.order.Type
}

func (e *OrderCanceled) ID() int64 {
	return e.order.ID
}

func (e *OrderCanceled) AccountID() int64 {
	return e.order.AccountID
}
