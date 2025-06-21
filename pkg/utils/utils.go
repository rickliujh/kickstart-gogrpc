package utils

import "log/slog"

func ParseSlogLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
