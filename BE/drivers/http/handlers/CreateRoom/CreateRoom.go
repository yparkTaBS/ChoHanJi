package CreateRoom

import (
	"ChoHanJi/infrastructure/Logging"
	RoomFactoryPorts "ChoHanJi/useCases/RoomFactory/ports"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CreateRoom struct {
	roomFactory RoomFactoryPorts.IRoomFactory
	validator   *validator.Validate
}

var _ http.Handler = (*CreateRoom)(nil)

func New(roomFactory RoomFactoryPorts.IRoomFactory, validator *validator.Validate) (*CreateRoom, error) {
	return &CreateRoom{roomFactory, validator}, nil
}

func (c *CreateRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	mapId, err := c.roomFactory.Create(data.MapWidth, data.MapHeight, data.Items)
	if err != nil {
		sendBack400(ctx, w, logger, "Failed to create room", err)
		return
	}

	res := Response{string(mapId)}
	responseBody, err := json.Marshal(res)
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
