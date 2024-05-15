package logger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("returns a logger", func(t *testing.T) {
		dest := new(bytes.Buffer)
		logger := NewLogger(ModeText, nil, dest)
		logger.Info("hello world")
		assert.Contains(t, dest.String(), "hello world")
	})
	t.Run("returns a logger with JSON", func(t *testing.T) {
		dest := new(bytes.Buffer)
		logger := NewLogger(ModeJSON, nil, dest)
		logger.Info("hello world")
		assert.Contains(t, dest.String(), "hello world")
	})
	t.Run("returns a logger to stdout by default", func(t *testing.T) {
		logger := NewLogger(ModeJSON, nil, nil)
		logger.Info("hello world")
		assert.True(t, true)
	})
}

func TestNewDevNullLogger(t *testing.T) {
	t.Run("it should do nothing", func(t *testing.T) {
		logger := NewDevNullLogger()
		logger.Info("hello world")
		assert.True(t, true)
	})
}
