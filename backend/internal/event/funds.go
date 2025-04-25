package event

type Funds struct {
	ID        int64
	AccountID int64
	Sum       float64
}

type FundsDebited struct {
	Funds Funds
}

type FundsCredited struct {
	Funds Funds
}

func (e *FundsDebited) Type() string {
	return FundsDebitedEvent
}

func (e *FundsDebited) ID() int64 {
	return e.Funds.ID
}

func (e *FundsDebited) AccountIDs() []int64 {
	return []int64{e.Funds.AccountID}
}

func (e *FundsDebited) Sum() float64 {
	return e.Funds.Sum
}

func (e *FundsCredited) Type() string {
	return FundsCreditedEvent
}

func (e *FundsCredited) ID() int64 {
	return e.Funds.ID
}

func (e *FundsCredited) AccountIDs() []int64 {
	return []int64{e.Funds.AccountID}
}

func (e *FundsCredited) Sum() float64 {
	return e.Funds.Sum
}
