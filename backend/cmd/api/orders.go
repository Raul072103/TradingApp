package main

import (
	"TradingSimulation/backend/internal/event"
	"go.uber.org/zap"
	"net/http"
)

type AddOrderRequest struct {
	AccountID int64  `json:"account_id"`
	StockID   int64  `json:"stock_id"`
	Quantity  int64  `json:"quantity"`
	OrderType string `json:"order_type"`
}

func (app *application) addOrderHandler(w http.ResponseWriter, r *http.Request) {
	var addOrderRequest AddOrderRequest
	err := readJSON(w, r, &addOrderRequest)
	if err != nil {
		app.logger.Error("error reading add order request", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	var placeOrder = event.OrderPlaced{
		AccountID: addOrderRequest.AccountID,
		Order: event.Order{
			AccountID: addOrderRequest.AccountID,
			Type:      addOrderRequest.OrderType,
			Count:     addOrderRequest.Quantity,
			Stock:     addOrderRequest.StockID,
		},
	}

	app.mainChannel <- &placeOrder

	err = writeJSON(w, http.StatusOK, struct {
		status string
	}{
		status: "order placed",
	})

	if err != nil {
		app.logger.Error("error sending JSON response", zap.Error(err))
	}
}
