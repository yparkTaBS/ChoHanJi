package Proceed

import (
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/ProceedUseCase"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type Struct struct {
	uc ProceedUseCase.Interface
}

var _ http.Handler = (*Struct)(nil)

func New(uc ProceedUseCase.Interface) *Struct {
	return &Struct{uc}
}

// ServeHTTP implements http.Handler.
func (s *Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := Logging.RetrieveLogger(ctx)
	if err != nil {
		sendBack500(ctx, w, logger, "Could not resolve the logger", err)
		return
	}

	// Get Data
	roomId := r.URL.Query().Get("roomId")
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.uc.Proceed(ctx, roomId, logger); err != nil {
		sendBack400(ctx, w, logger, "Something went wrong...", err)
		return
	}
}

func sendBack400(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, errMsg string, err error) {
	logger.ErrorContext(ctx, errMsg, slog.Any("Error", err))
	w.WriteHeader(http.StatusBadRequest)
}

func sendBack500(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, errMsg string, err error) {
	logger.ErrorContext(ctx, errMsg, slog.Any("Error", err))
	w.WriteHeader(http.StatusInternalServerError)
}
