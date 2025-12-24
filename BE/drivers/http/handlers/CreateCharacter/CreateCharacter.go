package CreateCharacter

import (
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/CharacterFactory"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CreateCharacter struct {
	uc        CharacterFactory.UseCaseInterface
	validator *validator.Validate
}

func New(uc CharacterFactory.UseCaseInterface, validator *validator.Validate) *CreateCharacter {
	return &CreateCharacter{uc, validator}
}

var _ http.Handler = (*CreateCharacter)(nil)

// ServeHTTP implements http.Handler.
func (c *CreateCharacter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if err := c.validator.Struct(data); err != nil {
		sendBack400(ctx, w, logger, "Request body failed at validation", err)
		return
	}

	characterId, err := c.uc.CreateCharacter(data.RoomId, data.UserName, data.Class, data.TeamNumber)
	if err != nil {
		sendBack400(ctx, w, logger, "Could not create the character", err)
		return
	}

	resp := Response{characterId}
	responseBody, err := json.Marshal(resp)
	if err != nil {
		sendBack500(ctx, w, logger, "CreateRoom.ServeHTTP: Could not marshal response", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write(responseBody); err != nil {
		sendBack500(ctx, w, logger, "CreateRoom.ServeHTTP: Failed to write response", err)
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
	RoomId     string `json:"RoomId" validate:"required"`
	UserName   string `json:"UserName" validate:"required"`
	Class      string `json:"Class" validate:"required"`
	TeamNumber int    `json:"TeamNumber" validate:"required,gt=0"`
}

type Response struct {
	CharacterId string `json:"CharacterId"`
}
