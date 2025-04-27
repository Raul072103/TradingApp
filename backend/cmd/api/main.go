package main

import (
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

	app := &application{
		logger: zapLogger,
		config: cfg,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	app.logger.Fatal("server error", zap.Error(app.run(mux)))

	err := app.run(mux)
	if err != nil {
		app.logger.Fatal("server error", zap.Error(err))
	}
}
