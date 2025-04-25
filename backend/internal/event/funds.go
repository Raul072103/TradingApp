package event

type FundDebited struct {
}

type FundsCredited struct {
}

func (e *FundDebited) Type() string {
	return FundsDebitedEvent
}

func (e *FundDebited) ID() int64 {
	return 0
}

func (e *FundDebited) AccountID() int64 {
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
