package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"nofelet/config"
	"nofelet/internal/app/chat"
	"nofelet/internal/dependency"

	"nofelet/pkg/httpserver"
)

func main() {
	if err := config.New(); err != nil {
		panic(err)
	}
	cfg := config.Current()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	deps, depErr := dependency.New(&cfg, logger)
	if depErr != nil {
		log.Fatal(depErr)
	}

	if sigErr := chat.New(deps); sigErr != nil {
		log.Fatal(sigErr)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	httpServer := httpserver.New(deps.Chat.Routes,
		httpserver.WithAddress(cfg.Chat.Port),
		httpserver.WithServerCRT(cfg.Crt),
		httpserver.WithServerKey(cfg.Key),
		httpserver.WithReadTimeout(cfg.Chat.ReadTimeout),
		httpserver.WithReadHeaderTimeout(cfg.Chat.ReadHeaderTimeout),
		httpserver.WithWriteTimeout(cfg.Chat.WriteTimeout),
		httpserver.WithShutdownTimeout(cfg.Chat.ShutdownTimeout),
	)

	select {
	case s := <-interrupt:
		logger.Error("error", slog.String("signal", s.String()))
	case err := <-httpServer.Notify():
		logger.Error("httpServer.Notify", slog.Any("error", err))
	}

	if err := httpServer.Shutdown(); err != nil {
		logger.Error("httpServer.Shutdown", slog.Any("error", err))
	}
}
