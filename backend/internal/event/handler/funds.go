package handler

type fundsHandler struct {
}

func (handler *fundsHandler) FundsCredited(sum float64) error {
	// extra funds credited logic here
	return nil
}

func (handler *fundsHandler) FundsDebited(sum float64) error {
	// extra funds debited logic here
	return nil
}
