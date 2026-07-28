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
	requestDuration metric.Float64Histogram
	activeRequests  metric.Int64UpDownCounter

	MessagesTotal              metric.Int64Counter
	UsersOnline                metric.Int64UpDownCounter
	RoomsActive                metric.Int64UpDownCounter
	WebSocketConnectionsActive metric.Int64UpDownCounter
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

	// Гистограмма времени выполнения запроса
	requestDuration, err = meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP requests."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		logger.Warn("сбой инициализации метрики requestDuration", err)
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

	// Количество активных пользователей
	UsersOnline, err = meter.Int64UpDownCounter(
		"chat.users.online",
		metric.WithDescription("Current number of online users."),
		metric.WithUnit("{user}"),
	)

	// Количество активных соединений
	WebSocketConnectionsActive, err = meter.Int64UpDownCounter(
		"chat.websocket.connections.active",
		metric.WithDescription("Current number of active WebSocket connections."),
		metric.WithUnit("{connection}"),
	)

	// Количество активных комнат
	RoomsActive, err = meter.Int64UpDownCounter(
		"chat.rooms.active",
		metric.WithDescription("Current number of active chat rooms."),
		metric.WithUnit("{room}"),
	)

	// Количество сообщений
	MessagesTotal, err = meter.Int64Counter(
		"chat.messages.total",
		metric.WithDescription("Total number of messages processed."),
		metric.WithUnit("{message}"),
	)
}
