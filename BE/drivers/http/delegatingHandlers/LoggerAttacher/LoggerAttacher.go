package LoggerAttacher

import (
	"ChoHanJi/infrastructure/ContextKeys"
	"ChoHanJi/infrastructure/Logging"
	"context"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type LoggerAttacher struct {
	next http.Handler
}

var _ http.Handler = (*LoggerAttacher)(nil)

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return new(next)
	}
}

func new(next http.Handler) *LoggerAttacher {
	return &LoggerAttacher{next}
}

func (l *LoggerAttacher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.Default()
	logger = logger.WithGroup("Request")

	requestId, err := uuid.NewV7()
	if err != nil {
		logger.ErrorContext(ctx, "Error Generating RequestId for some reason...", slog.Any("Error", err))
	}

	jobName, err := Logging.GetJobName(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Error getting the JobName", slog.Any("Error", err))
	}

	logger = logger.With("RequestID", requestId)
	logger = logger.With("JobName", jobName)

	ctx = context.WithValue(ctx, ContextKeys.Logger, logger)
	r = r.WithContext(ctx)

	l.next.ServeHTTP(w, r)
}
