package StartGame

import (
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/StartGameUseCase"
	"log/slog"
	"net/http"
	"strings"
)

type StartGame struct {
	uc StartGameUseCase.Interface
}

var _ http.Handler = (*StartGame)(nil)

func New(uc StartGameUseCase.Interface) *StartGame {
	return &StartGame{uc}
}

// ServeHTTP implements http.Handler.
func (s *StartGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, _ := Logging.RetrieveLogger(ctx)

	roomId := r.URL.Query().Get("roomId")
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.uc.Announce(roomId); err != nil {
		logger.ErrorContext(ctx, "Something went wrong...", slog.Any("Error", err))
	}
}
