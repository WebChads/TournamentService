package slogerr

import "log/slog"

func Error(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Warn(err error) slog.Attr {
	return slog.Attr{
		Key: "warning",
		Value: slog.StringValue(err.Error()),
	}
}