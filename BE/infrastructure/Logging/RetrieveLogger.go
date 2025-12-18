package Logging

import (
	"ChoHanJi/infrastructure/ContextKeys"
	"context"
	"errors"
	"log/slog"
)

var (
	ErrLoggerNotFound    = errors.New("logger not found in the context")
	ErrLoggerTypeUnknown = errors.New("logger is not of slog.Logger")
)

func RetrieveLogger(ctx context.Context) (*slog.Logger, error) {
	registration := ctx.Value(ContextKeys.Logger)
	if registration == nil {
		return &slog.Logger{}, ErrLoggerNotFound
	}

	logger, ok := registration.(*slog.Logger)
	if !ok {
		return &slog.Logger{}, ErrLoggerTypeUnknown
	}

	return logger, nil
}
