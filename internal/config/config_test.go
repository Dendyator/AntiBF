package config

import (
	"testing"

	"github.com/Dendyator/AntiBF/internal/logger" //nolint
)

func TestLoadConfig(t *testing.T) {
	logLevel := "info"
	appLogger := logger.New(logLevel)
	config := LoadConfig("../../configs/config.yaml", appLogger)

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected 0.0.0.0, got %s", config.Server.Host)
	}
}
