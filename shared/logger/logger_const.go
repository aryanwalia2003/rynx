package logger

import (
	"log/slog"
	"os"
)

func Logger_const() *Logger_struct {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	instance := slog.New(handler)
	
	return &Logger_struct{
		instance: instance,
	}
}
