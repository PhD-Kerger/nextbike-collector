package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/kjk/common/filerotate"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if e := h.Handle(ctx, r); e != nil && err == nil {
				err = e
			}
		}
	}
	return err
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: hs}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: hs}
}

var Logger *slog.Logger = slog.Default()

// SetupLogger initializes the global logger with multi-handler (stdout/info and file/warn)
// Uses filerotate for log rotation
func SetupLogger(logFilePath string) error {
	// Configure filerotate for log rotation
	logFile, err := filerotate.NewDaily(logFilePath, "", nil)
	if err != nil {
		return err
	}

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo})
	multiHandler := &MultiHandler{handlers: []slog.Handler{stdoutHandler, fileHandler}}
	Logger = slog.New(multiHandler)
	return nil
}
