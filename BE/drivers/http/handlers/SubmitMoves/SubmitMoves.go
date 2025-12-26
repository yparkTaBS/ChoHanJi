package SubmitMoves

import (
	"ChoHanJi/domain/MoveFlag"
	"ChoHanJi/domain/Player"
	"ChoHanJi/infrastructure/Logging"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Struct struct {
	validator *validator.Validate
}

var _ http.Handler = (*Struct)(nil)

func New(validator *validator.Validate) *Struct {
	return &Struct{validator}
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

	// Unmarshal
	var data Request
	if err := json.Unmarshal(request, &data); err != nil {
		sendBack500(ctx, w, logger, "Could not unmarshal the body", err)
		return
	}

	// Validate
	if err := s.validator.Struct(data); err != nil {
		sendBack400(ctx, w, logger, "Request body failed at validation", err)
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

type Request struct {
	ActionType MoveFlag.Enum `json:"ActionType" validate:"required"`
	X          int           `json:"X" validate:"required,gte=0"`
	Y          int           `json:"Y" validate:"required,gte=0"`
	SubjectStr Player.Id     `json:"Subject" validate:"required,alphanum,len=5"`
}
