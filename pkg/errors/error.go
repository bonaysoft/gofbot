package errors

import "log/slog"

func With(err error) slog.Attr {
	return slog.Any("error", err)
}
