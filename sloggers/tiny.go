package sloggers

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type TinyHandler struct {
	slog.Handler
	l *log.Logger
}

func NewTinyHandler(out io.Writer, opts *slog.HandlerOptions) *TinyHandler {
	return &TinyHandler{
		Handler: slog.NewTextHandler(out, opts),
		l:       log.New(out, "", 0),
	}
}

func (h *TinyHandler) Handle(ctx context.Context, r slog.Record) error {
	var level string

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString("DBG")
	case slog.LevelInfo:
		level = color.BlueString("INF")
	case slog.LevelWarn:
		level = color.YellowString("WRN")
	case slog.LevelError:
		level = color.RedString("ERR")
	}

	attrs := make([]string, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		s := color.BlackString("%s: %s", a.Key, a.Value.String())
		attrs = append(attrs, s)

		return true
	})

	lines := []string{fmt.Sprintf("%s %s", level, r.Message)}
	lines = append(lines, attrs...)
	sep := color.HiBlackString("\n└─ ")

	h.l.Println(strings.Join(lines, sep))

	return nil
}
