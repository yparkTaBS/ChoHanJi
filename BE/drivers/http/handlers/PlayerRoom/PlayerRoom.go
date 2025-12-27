package PlayerRoom

import (
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/PlayerWaitingRoomUseCase"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type PlayerRoom struct {
	uc PlayerWaitingRoomUseCase.UseCaseInterface
}

func New(uc PlayerWaitingRoomUseCase.UseCaseInterface) *PlayerRoom {
	return &PlayerRoom{uc}
}

var _ http.Handler = (*PlayerRoom)(nil)

func (p *PlayerRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	logger, _ := Logging.RetrieveLogger(ctx)

	ioWriter, ok := w.(io.Writer)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	playerId := r.URL.Query().Get("playerId")
	if len(strings.TrimSpace(roomId)) == 0 || len(strings.TrimSpace(playerId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flusher.Flush()

	if err := p.uc.ConnectAndListen(ctx, ioWriter, roomId, playerId, flusher); err != nil {
		logger.ErrorContext(ctx, "Something went wrong...", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
