package logger

import "testing"

func TestLogger_const(t *testing.T) {
	log := Logger_const()
	if log == nil {
		t.Fatal("expected logger to be initialized")
	}

	if log.instance == nil {
		t.Fatal("expected internal slog instance to be initialized")
	}
}
