package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shoter/internal/config"
	"url-shoter/internal/http-sever/handlers/url/save"
	mvLogger "url-shoter/internal/http-sever/middleware/logger"
	"url-shoter/internal/lib/logger/handlers/slogpretty"
	"url-shoter/internal/lib/logger/sl"
	"url-shoter/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// init config clean env
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	//init logger  slog
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	//init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = fmt.Sprint(storage)

	// init router
	router := chi.NewRouter()
	//middleware
	router.Use(middleware.RequestID)
	router.Use(mvLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	//todo run server
	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.HttpServer.Timeout,
		WriteTimeout:      cfg.HttpServer.Timeout,
		IdleTimeout:       cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed start server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
