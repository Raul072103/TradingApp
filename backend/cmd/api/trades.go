package main

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (app *application) executedTradesHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare a list of simplified trade info
	type TradeInfo struct {
		TradeID    int64     `json:"trade_id"`
		AccountIDs []int64   `json:"account_ids"`
		Timestamp  time.Time `json:"timestamp"`
	}

	var executedTrades []TradeInfo

	for _, event := range app.materializedView.Trades {
		executedTrades = append(executedTrades, TradeInfo{
			TradeID:    event.Trade.ID,
			AccountIDs: event.Trade.AccountIDs,
			Timestamp:  event.Timestamp,
		})
	}

	err := writeJSON(w, http.StatusOK, executedTrades)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		app.logger.Error("error retrieving executed trades", zap.Error(err))
	}
}
