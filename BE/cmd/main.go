package main

import (
	"ChoHanJi/CompositionRoot"
	"ChoHanJi/bootstrap"
	"ChoHanJi/config/PilgrimCraftConfig"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
)

func main() {
	appContext := context.Background()

	config := PilgrimCraftConfig.LoadSettings(appContext)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.MinimumLoggingLevel,
	}))
	slog.SetDefault(logger)

	container := bootstrap.Register(appContext)

	routes, err := CompositionRoot.CreateEndPoints(container, config)
	if err != nil {
		panic(err)
	}

	serverAddr := fmt.Sprintf(":%s", config.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: routes,
		BaseContext: func(net.Listener) context.Context {
			return appContext
		},
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.ErrorContext(appContext, fmt.Errorf("Could not listen on 8080: %w", err).Error())
	}
}
