package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"nofelet/config"
	"nofelet/internal/app/chat"
	"nofelet/internal/dependency"
	"nofelet/middleware/metrics"

	"nofelet/pkg/httpserver"
)

func main() {
	if err := config.New(); err != nil {
		panic(err)
	}
	cfg := config.Current()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Создаем современный экспортер v0.66.0
	exporter, err := otelprom.New()
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(provider)
	metrics.Init()
	defer func() { _ = provider.Shutdown(context.Background()) }()

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
