package event

type FundsDebited struct {
}

type FundsCredited struct {
}

func (e *FundsDebited) Type() string {
	return FundsDebitedEvent
}

func (e *FundsDebited) ID() int64 {
	return 0
}

func (e *FundsDebited) AccountID() int64 {
	return 0
}

func (e *FundsCredited) Type() string {
	return FundsCreditedEvent
}

func (e *FundsCredited) ID() int64 {
	return 0
}

func (e *FundsCredited) AccountID() int64 {
	return 0
}
