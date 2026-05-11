package config

import (
	"fmt"
)

var current Config

// New - создает объект конфига
func New() error {
	c, err := newConfig()
	if err != nil {
		return fmt.Errorf("fail to create new config %w", err)
	}
	current = *c

	return nil
}

func Current() Config {
	return current
}
