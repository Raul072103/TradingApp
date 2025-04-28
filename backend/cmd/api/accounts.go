package main

import (
	"TradingSimulation/backend/internal/event"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (app *application) accountOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var accountOrders = make(map[int64]struct {
		BuyOrders        []event.Order `json:"buy_orders"`
		SellOrders       []event.Order `json:"sell_orders"`
		SuccessfulOrders []event.Order `json:"successful_orders"`
		CanceledOrders   []event.Order `json:"canceled_orders"`
	})

	for key := range app.materializedView.Accounts {
		account := app.materializedView.Accounts[key]

		var canceledOrders = make([]event.Order, 0)

		for key := range account.CanceledOrders {
			canceledOrders = append(canceledOrders, account.CanceledOrders[key])
		}

		accountOrders[key] = struct {
			BuyOrders        []event.Order `json:"buy_orders"`
			SellOrders       []event.Order `json:"sell_orders"`
			SuccessfulOrders []event.Order `json:"successful_orders"`
			CanceledOrders   []event.Order `json:"canceled_orders"`
		}{
			BuyOrders:        account.BuyOrders,
			SellOrders:       account.SellOrders,
			SuccessfulOrders: account.SuccessfulOrders,
			CanceledOrders:   canceledOrders}
	}

	err := writeJSON(w, http.StatusOK, accountOrders)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		app.logger.Error("error getting all accounts", zap.Error(err))
	}
}

func (app *application) accountsFundsHandler(w http.ResponseWriter, r *http.Request) {
	accountFunds := make(map[int64]float64)

	for accountID, accountState := range app.materializedView.Accounts {
		accountFunds[accountID] = accountState.Funds
	}

	err := writeJSON(w, http.StatusOK, accountFunds)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		app.logger.Error("error getting account funds", zap.Error(err))
	}
}

func (app *application) accountOrdersByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get the account ID from the URL parameters
	accountIDParam := chi.URLParam(r, "id")

	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Find the account in the materialized view
	account, exists := app.materializedView.Accounts[int64(accountID)]
	if !exists {
		err := writeJSONError(w, http.StatusNotFound, "Account not found")
		if err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		return
	}

	// Prepare the response structure with Buy and Sell orders
	accountOrders := struct {
		BuyOrders  []event.Order `json:"buy_orders"`
		SellOrders []event.Order `json:"sell_orders"`
	}{
		BuyOrders:  account.BuyOrders,
		SellOrders: account.SellOrders,
	}

	// Write the response
	err = writeJSON(w, http.StatusOK, accountOrders)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		app.logger.Error("error retrieving account orders", zap.Error(err))
	}
}
