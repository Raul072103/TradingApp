package main

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/handler"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
	"TradingSimulation/backend/internal/filter"
	"TradingSimulation/common/logger"
	"expvar"
	"go.uber.org/zap"
	"runtime"
)

const (
	version = "0.0.0"
)

func main() {
	zapLogger := logger.InitLogger("backend.log")

	cfg := config{
		addr: "localhost:8080",
	}

	mainChannel := make(chan event.Event, 10000)
	processedEvents := make(chan event.Event, 10000)

	eventStore, err := store.New()
	if err != nil {
		zapLogger.Fatal("server error", zap.Error(err))
	}
	defer eventStore.Close()

	existingEvents, err := eventStore.GetAllEvents()
	if err != nil {
		zapLogger.Fatal("server error", zap.Error(err))
	}

	materializedView, err := view.New(existingEvents)
	if err != nil {
		zapLogger.Fatal("server error", zap.Error(err))
	}

	mainHandler := handler.New(mainChannel, processedEvents, &materializedView, materializedView.Stocks)

	processedEventsFilter := filter.New(eventStore, &materializedView, processedEvents)

	app := &application{
		logger:                zapLogger,
		config:                cfg,
		mainHandler:           &mainHandler,
		materializedView:      &materializedView,
		eventStore:            eventStore,
		processedEventsFilter: processedEventsFilter,
		mainChannel:           mainChannel,
		processedEvents:       processedEvents,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	app.logger.Fatal("server error", zap.Error(app.run(mux)))

	err = app.run(mux)
	if err != nil {
		app.logger.Fatal("server error", zap.Error(err))
	}
}
