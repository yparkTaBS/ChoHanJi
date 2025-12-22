package WaitingRoom

import (
	"ChoHanJi/useCases/AdminWaitingRoomUseCase"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type AdminWaitingRoom struct {
	uc AdminWaitingRoomUseCase.IAdminWaitingRoomUseCase
}

var _ http.Handler = (*AdminWaitingRoom)(nil)

func New(useCase AdminWaitingRoomUseCase.IAdminWaitingRoomUseCase) *AdminWaitingRoom {
	return &AdminWaitingRoom{useCase}
}

func (a *AdminWaitingRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

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
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	if err := a.uc.ConnectAndListen(ctx, ioWriter, roomId, flusher); err != nil {
		slog.Default().ErrorContext(ctx, "Something went wrong...", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
