package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerInitialization(t *testing.T) {
	var buf bytes.Buffer

	log := New("debug")
	log.Out = &buf

	assert.Equal(t, logrus.DebugLevel, log.Level, "Expected log level to be debug")

	log.Debug("Test message")

	assert.Contains(t, buf.String(), "Test message", "Expected log message to be present")
}

func TestLoggerInitializationWithInvalidLevel(t *testing.T) {
	var buf bytes.Buffer

	log := New("invalid-level")
	log.Out = &buf

	assert.Equal(t, logrus.InfoLevel, log.Level, "Expected log level to default to info")

	log.Info("Test message with invalid level")

	assert.Contains(t, buf.String(), "Test message with invalid level", "Expected log message to be present with default info level")
}
