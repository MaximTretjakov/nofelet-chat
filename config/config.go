package config

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	WS     WSConfig `env:",prefix=WS_"`     // Инфа сервера
	CoTURN CoTURN   `env:",prefix=COTURN_"` // Котерн инфа
	Debug  bool     `env:"DEBUG"`           // Дебаг режимы
}

type WSConfig struct {
	Port              string        `env:"PORT,required"`                   // Порт
	ServerCrt         string        `env:"SERVER_CRT,required"`             // Сертификат
	ServerKey         string        `env:"SERVER_KEY,required"`             // Сертификат
	ReadTimeout       time.Duration `env:"READ_TIMEOUT,default=30s"`        // Таймаут на чтение
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT,default=30s"`       // Таймаут на запись
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT,default=30s"` // Таймаут на чтение хедеров
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,default=3s"`     // Таймаут на завершение
}

type CoTURN struct {
	SharedSecret string `env:"SHARED_SECRET,required"`  // Секрет
	TurnServerIP string `env:"TURN_SERVER_IP,required"` // IP адресе сервера
}

func init() {
	_ = godotenv.Load()
}

// NewConfig - выгружает данные из .env и осздает переменные окружения
func newConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process(context.Background(), &config); err != nil {
		return nil, err
	}

	return &config, nil
}
