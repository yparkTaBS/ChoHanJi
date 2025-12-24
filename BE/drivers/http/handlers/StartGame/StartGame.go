package StartGame

import (
	"net/http"
	"strings"
)

type StartGame struct{}

var _ http.Handler = (*StartGame)(nil)

func New() *StartGame {
	return &StartGame{}
}

// ServeHTTP implements http.Handler.
func (s *StartGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	if len(strings.TrimSpace(roomId)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
