package main

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func (app *application) goroutineCountHandler(w http.ResponseWriter, r *http.Request) {
	numGoroutines := runtime.NumGoroutine()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return the number of goroutines in JSON format
	response := map[string]int{"goroutines": numGoroutines}
	json.NewEncoder(w).Encode(response)
}
