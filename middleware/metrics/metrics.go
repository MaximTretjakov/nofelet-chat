package metrics

import (
	"github.com/bytedance/gopkg/util/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const serviceLabel = "nofelet-chat"

var (
	requestsTotal   metric.Int64Counter
	responseTotal   metric.Int64Counter
	errorsTotal     metric.Int64Counter
	panicsTotal     metric.Int64Counter
	responseSize    metric.Int64Histogram
	requestDuration metric.Float64Histogram
	activeRequests  metric.Int64UpDownCounter
)

func Init() {
	var err error

	// Создаем Meter (пространство имен для метрик вашего приложения)
	meter := otel.Meter(serviceLabel)

	// Счетчик всех HTTP-запросов
	requestsTotal, err = meter.Int64Counter(
		"http.server.requests.total",
		metric.WithDescription("Total number of HTTP requests."),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики requestsTotal", err)
	}

	// Счетчик всех HTTP-ответов
	responseTotal, err = meter.Int64Counter(
		"http.server.responses.total",
		metric.WithDescription("Total number of HTTP responses."),
		metric.WithUnit("{response}"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики responseTotal", err)
	}

	// Счетчик ошибок
	errorsTotal, err = meter.Int64Counter(
		"http.server.errors.total",
		metric.WithDescription("Total number of failed HTTP requests."),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики errorsTotal", err)
	}

	// Счетчик паник
	panicsTotal, err = meter.Int64Counter(
		"http.server.panics.total",
		metric.WithDescription("Total number of recovered panics."),
		metric.WithUnit("{panic}"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики panicsTotal", err)
	}

	// Гистограмма времени выполнения запроса
	requestDuration, err = meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP requests."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики requestDuration", err)
	}

	// Размер ответа
	responseSize, err = meter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("HTTP response size."),
		metric.WithUnit("By"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики responseSize", err)
	}

	// Активные запросы
	activeRequests, err = meter.Int64UpDownCounter(
		"http.server.requests.active",
		metric.WithDescription("Current number of active HTTP requests."),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики activeRequests", err)
	}
}
