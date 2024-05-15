package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	ModeText = "text"
	ModeJSON = "json"
)

func NewLogger(mode string, opts *slog.HandlerOptions, destination io.Writer) *slog.Logger {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if destination == nil {
		destination = os.Stdout
	}
	opts.AddSource = true
	if strings.ToLower(mode) == ModeJSON {
		return slog.New(slog.NewJSONHandler(destination, opts))
	}
	return slog.New(slog.NewTextHandler(destination, opts))
}

type DevNull struct{}

// Write implements the io.Writer interface but does nothing
func (d *DevNull) Write(p []byte) (n int, err error) { return }

func NewDevNullLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&DevNull{}, &slog.HandlerOptions{}))
}
