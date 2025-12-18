package CompositionRoot

import (
	"ChoHanJi/drivers/http/delegatingHandlers/GenericPanicCatcher"
	"ChoHanJi/drivers/http/delegatingHandlers/JobNameAttacher"
	"ChoHanJi/drivers/http/delegatingHandlers/LoggerAttacher"
	"ChoHanJi/drivers/http/handlers"
	"ChoHanJi/infrastructure/Logging"
	ctx "context"
	"log/slog"
	"net/http"
	"time"

	"github.com/TaBSRest/GoFac"
	gi "github.com/TaBSRest/GoFac/interfaces"
	"github.com/go-chi/chi"
)

func CreateEndPoints(container gi.Container) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Mount("/api", RegisterPOSTRoom(container))

	return r, nil
}

func RegisterPOSTRoom(container gi.Container) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New("POST /api/room"))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Post("/room", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		context, cancel := ctx.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()

		r = r.WithContext(context)

		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "POST /api/room: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(handlers.POSTRoom))
		if err != nil {
			logger.ErrorContext(context, "POST /api/room: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}
