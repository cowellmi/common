package env

import (
	"log/slog"
	"os"

	"github.com/cowellmi/common/env/parse"
)

var logger = slog.Default()

func SetLogger(l *slog.Logger) {
	logger = l
}

func Get[T any](key string, defaultValue T, parse parse.Parser[T]) T {
	s, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	v, err := parse(s)
	if err != nil {
		logger.Warn(
			"failed to parse environment variable value",
			slog.String("key", key),
			slog.String("value", s),
			slog.String("error", err.Error()),
			slog.Any("default", defaultValue),
		)
		return defaultValue
	}

	return v
}
