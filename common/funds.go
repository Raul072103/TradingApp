package common

type Funds struct {
	ID        int64
	AccountID int64
	Sum       float64
}

type FundsDebited struct {
	Funds Funds
	Trade Trade
}

type FundsCredited struct {
	Funds Funds
	Trade Trade
}

func (e *FundsDebited) Type() int {
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

func (e *FundsCredited) Type() int {
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
