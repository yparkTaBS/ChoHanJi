package SubmitFightResult

import (
	"ChoHanJi/domain/Fight"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/SubmitFightResultUseCase"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Struct struct {
	uc        SubmitFightResultUseCase.Interface
	validator *validator.Validate
}

var _ http.Handler = (*Struct)(nil)

func New(uc SubmitFightResultUseCase.Interface, validator *validator.Validate) *Struct {
	return &Struct{uc, validator}
}

func (s *Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := Logging.RetrieveLogger(ctx)
	if err != nil {
		sendBack500(ctx, w, logger, "Could not resolve the logger", err)
		return
	}

	body := r.Body
	defer body.Close()

	requestBytes, err := io.ReadAll(body)
	if err != nil {
		sendBack500(ctx, w, logger, "Could not read the request", err)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req SubmitFightResultUseCase.Request
	if err := json.Unmarshal(requestBytes, &req); err != nil {
		sendBack400(ctx, w, logger, "Wrong Submission", err)
		return
	}

	if err := s.validator.Struct(req); err != nil {
		sendBack400(ctx, w, logger, "Wrong Submission", err)
		return
	}

	if err := s.uc.Submit(Room.Id(roomId), Fight.Id(req.FightId), Player.Id(req.SubmitterId), Player.Id(req.WinnerId)); err != nil {
		sendBack500(ctx, w, logger, "Something went wrong...", err)
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
