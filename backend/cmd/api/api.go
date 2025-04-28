package main

import (
	"TradingSimulation/backend/internal/event"
	"TradingSimulation/backend/internal/event/handler"
	"TradingSimulation/backend/internal/event/store"
	"TradingSimulation/backend/internal/event/view"
	"TradingSimulation/backend/internal/filter"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO() don't forget about filter!

type application struct {
	config                config
	logger                *zap.Logger
	eventStore            *store.Store
	materializedView      *view.MaterializedView
	mainHandler           *handler.Handler
	processedEventsFilter *filter.Filter
	mainChannel           chan event.Event
	processedEvents       chan event.Event
}

type config struct {
	addr        string
	env         string
	db          dbConfig
	apiURL      string
	frontendURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Get("/executed_trades", app.executedTradesHandler)

		r.Get("/accounts", app.accountsFundsHandler)
		r.Get("/accounts/orders", app.accountOrdersHandler)
		r.Get("/accounts/{id}/orders", app.accountOrdersByIDHandler)

		r.Post("/orders", app.addOrderHandler)

		r.Get("/metrics/goroutines", app.goroutineCountHandler)
	})

	return mux
}

func (app *application) run(mux *chi.Mux) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// graceful shutdown
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("signal caught", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	go func() {
		err := app.processedEventsFilter.Run()
		if err != nil {
			shutdown <- err
		}
	}()

	go func() {
		err := app.mainHandler.Run()
		if err != nil {
			shutdown <- err
		}
	}()

	app.logger.Info("Server has started", zap.String("addr", app.config.addr), zap.String("env", app.config.env))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Info("server has stopped", zap.String("addr", app.config.addr), zap.String("env", app.config.env))

	return nil
}
