package GenericPanicCatcher

import (
	"ChoHanJi/infrastructure/Logging"
	ctx "context"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func New() func(http.Handler) http.Handler {
	return new
}

func new(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := r.Context()
		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			panic(err)
		}

		jobName, err := Logging.GetJobName(context)
		if err != nil {
			panic(err)
		}

		defer func() {
			if rec := recover(); rec != nil {
				logPanic(logger, context, rec, jobName)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func logPanic(logger *slog.Logger, context ctx.Context, rec interface{}, location string) {
	locationAttribute := slog.Attr{Key: "location", Value: slog.StringValue(location)}
	stackTrace := slog.Attr{Key: "Stack Trace", Value: slog.StringValue(string(debug.Stack()))}
	switch typedRecover := rec.(type) {
	case string:
		logger.ErrorContext(context, "Internal Server Error", locationAttribute, slog.String("error", typedRecover), stackTrace)
	case error:
		logger.ErrorContext(context, "Internal Server Error", locationAttribute, slog.Any("error", typedRecover), stackTrace)
	default:
		logger.ErrorContext(context, "Internal Server Error", locationAttribute, slog.Any("error", typedRecover), stackTrace)
	}
}
