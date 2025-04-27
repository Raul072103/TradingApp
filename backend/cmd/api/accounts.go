package main

import (
	"go.uber.org/zap"
	"net/http"
)

func (app *application) accountTradesHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) accountOrdersHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) allAccountsHandler(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(w, http.StatusOK, app.materializedView.Accounts)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			app.logger.Error("error writing error response", zap.Error(err))
		}
		app.logger.Error("error getting all accounts", zap.Error(err))
	}
}
