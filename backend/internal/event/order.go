package event

type OrderPlaced struct {
}

type OrderCanceled struct {
}

func (e *OrderPlaced) Type() string {
	return OrdersPlacedEvent
}

func (e *OrderPlaced) ID() int64 {
	return 0
}

func (e *OrderPlaced) AccountID() int64 {
	return 0
}

func (e *OrderCanceled) Type() string {
	return OrdersCanceledEvent
}

func (e *OrderCanceled) ID() int64 {
	return 0
}

func (e *OrderCanceled) AccountID() int64 {
	return 0
}
