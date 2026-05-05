package logger

import "log/slog"

func (l *Logger_struct) Logger_error_method(message string, err error) {
	l.instance.Error(message, slog.Any("error", err))
}
