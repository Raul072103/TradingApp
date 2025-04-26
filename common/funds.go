package common

type Funds struct {
	ID        int64   `json:"id"`
	AccountID int64   `json:"accountID"`
	Sum       float64 `json:"sum"`
}

type FundsDebited struct {
	Funds Funds `json:"funds"`
	Trade Trade `json:"trade"`
}

type FundsCredited struct {
	Funds Funds `json:"funds"`
	Trade Trade `json:"trade"`
}

func (e *FundsDebited) Type() int64 {
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

func (e *FundsCredited) Type() int64 {
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
