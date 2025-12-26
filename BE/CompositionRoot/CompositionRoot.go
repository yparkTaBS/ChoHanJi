package CompositionRoot

import (
	"ChoHanJi/config/PilgrimCraftConfig"
	"ChoHanJi/drivers/http/delegatingHandlers/GenericPanicCatcher"
	"ChoHanJi/drivers/http/delegatingHandlers/JobNameAttacher"
	"ChoHanJi/drivers/http/delegatingHandlers/LoggerAttacher"
	"ChoHanJi/drivers/http/handlers"
	"ChoHanJi/infrastructure/Logging"
	ctx "context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/TaBSRest/GoFac"
	gi "github.com/TaBSRest/GoFac/interfaces"
	"github.com/go-chi/chi"
)

const (
	GET     = "GET"
	POST    = "POST"
	OPTIONS = "OPTIONS"
)

func CreateEndPoints(container gi.Container, config *PilgrimCraftConfig.PilgrimCraftConfig) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Default().ErrorContext(r.Context(), fmt.Sprintf("Not found %s", r.URL.Path))
		w.WriteHeader(http.StatusNotFound)
	}))

	origin := fmt.Sprintf("%s:%s", config.Server.Host, "3000")

	r.Mount("/api/room", RegisterPOSTRoom(container, handlers.POSTRoom, origin))
	r.Mount("/api/character", RegisterPOSTPlayer(container, handlers.POSTCharacter, origin))
	r.Mount(string(handlers.GETPlayerEvent), RegisterGETPlayerEvent(container, string(handlers.GETPlayerEvent), origin))
	r.Mount("/api/room/waiting/admin", RegisterAdminWaitingRoom(container, handlers.GETRoomAdmin, origin))
	r.Mount(string(handlers.POSTGameStart), RegisterPOSTGameStart(container, handlers.POSTGameStart, origin))

	r.Mount(string(handlers.GETAdminGameStatus), RegisterEndPoint(container, GET, string(handlers.GETAdminGameStatus), origin))
	r.Mount(string(handlers.GETPlayerGameStatus), RegisterEndPoint(container, GET, string(handlers.GETPlayerGameStatus), origin))

	r.Mount(string(handlers.POSTSubmitMoves), RegisterEndPoint(container, POST, string(handlers.POSTSubmitMoves), origin))

	return r, nil
}

func RegisterPOSTRoom(container gi.Container, route handlers.RouteToken, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("POST %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		context, cancel := ctx.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()

		r = r.WithContext(context)

		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "POST /api/room: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "POST /api/room: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterPOSTPlayer(container gi.Container, route handlers.RouteToken, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("POST %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		context, cancel := ctx.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()

		r = r.WithContext(context)

		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "POST /api/character: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "POST /api/character: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterPOSTGameStart(container gi.Container, route handlers.RouteToken, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("POST %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		context, cancel := ctx.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()

		r = r.WithContext(context)

		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "POST /api/game/start: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "POST /api/game/start: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterAdminWaitingRoom(container gi.Container, route handlers.RouteToken, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("GET %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		context := r.Context()
		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "GET /api/room/admin: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "GET /api/room/admin: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterGETPlayerEvent(container gi.Container, route string, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("GET %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		context := r.Context()
		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "GET /api/player/event: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "POST /api/player/event: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterGETGameStatus(container gi.Container, route, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("GET %s", route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		context := r.Context()
		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, "GET /api/game: Could not retrieve logger from the context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, "GET /api/game: Could not resolve handler", slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}

func RegisterEndPoint(container gi.Container, verb, route, origin string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JobNameAttacher.New(fmt.Sprintf("%s %s", verb, route)))
	r.Use(LoggerAttacher.New())
	r.Use(GenericPanicCatcher.New())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("%s, %s", verb, OPTIONS))

		context := r.Context()
		logger, err := Logging.RetrieveLogger(context)
		if err != nil {
			slog.ErrorContext(context, fmt.Sprintf("%s %s: Could not retrieve logger from the context", verb, route))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler, err := GoFac.ResolveNamed[http.Handler](container, context, string(route))
		if err != nil {
			logger.ErrorContext(context, fmt.Sprintf("%s %s: Could not resolve handler", verb, route), slog.Any("Error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})

	return r
}
