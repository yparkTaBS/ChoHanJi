package SkipMove

import (
	"ChoHanJi/domain/Action"
	"ChoHanJi/domain/Room"
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/SubmitMoveUseCase"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Struct struct {
	uc SubmitMoveUseCase.Interface
}

var _ http.Handler = (*Struct)(nil)

func New(uc SubmitMoveUseCase.Interface) *Struct {
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
	requestBody := r.Body
	defer requestBody.Close()

	request, err := io.ReadAll(requestBody)
	if err != nil {
		sendBack500(ctx, w, logger, "Could not read the request", err)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.uc.Submit(Room.Id(roomId), Action.Skip, request); err != nil {
		switch {
		case errors.Is(err, SubmitMoveUseCase.ErrWrongInput):
			sendBack400(ctx, w, logger, "Wrong Submission", err)
		default:
			sendBack500(ctx, w, logger, "Something went wrong...", err)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func sendBack400(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, errMsg string, err error) {
	logger.ErrorContext(ctx, errMsg, slog.Any("Error", err))
	w.WriteHeader(http.StatusBadRequest)
}

func sendBack500(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, errMsg string, err error) {
	logger.ErrorContext(ctx, errMsg, slog.Any("Error", err))
	w.WriteHeader(http.StatusInternalServerError)
}
