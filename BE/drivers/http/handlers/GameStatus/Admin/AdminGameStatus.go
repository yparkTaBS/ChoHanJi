package AdminGameStatus

import (
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/GameStatus"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type AdminGameStatus struct {
	uc GameStatus.Interface
}

var _ http.Handler = (*AdminGameStatus)(nil)

// ServeHTTP implements http.Handler.
func (a *AdminGameStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	if err := a.uc.ConnectAndListen(ctx, ioWriter, roomId, "admin", flusher); err != nil {
		logger.ErrorContext(ctx, "Something went wrong...", slog.Any("Error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
