package decorator

import (
	"log/slog"

	"github.com/gorilla/websocket"

	"nofelet/config"
)

// ConnectionManager - оборачивает стандартный websocket.Conn
type ConnectionManager struct {
	*websocket.Conn
	cfg      *config.Config
	logger   *slog.Logger
	remoteIP string
}

// NewConn - создает логирующий декоратор
func NewConn(conn *websocket.Conn, logger *slog.Logger, cfg *config.Config) *ConnectionManager {
	return &ConnectionManager{
		Conn:     conn,
		cfg:      cfg,
		logger:   logger,
		remoteIP: conn.RemoteAddr().String(),
	}
}

// WriteJSON - логирует отправку данных
func (cm *ConnectionManager) WriteJSON(v interface{}) error {
	err := cm.Conn.WriteJSON(v)

	if cm.cfg.Debug {
		cm.logger.Info("ws_msg_sent",
			slog.String("to", cm.remoteIP),
			slog.Any("data", v),
		)
	}

	return err
}

// ReadJSON - перехватывает получение данных
func (cm *ConnectionManager) ReadJSON(v interface{}) error {
	err := cm.Conn.ReadJSON(v)
	if err != nil {
		return err
	}

	if cm.cfg.Debug {
		cm.logger.Info("ws_msg_received",
			slog.String("from", cm.remoteIP),
			slog.Any("data", v),
		)
	}

	return err
}
