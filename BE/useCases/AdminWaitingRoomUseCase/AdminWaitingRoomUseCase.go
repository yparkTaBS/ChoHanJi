package AdminWaitingRoomUseCase

import (
	r "ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"ChoHanJi/infrastructure/Logging"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type IHub interface {
	Subscribe(roomId, subscriberId string) <-chan []byte
	Unsubscribe(roomId, subscriberId string) error
}

var _ IHub = (*SSEHub.SSEHub)(nil)

type UseCaseInterface interface {
	ConnectAndListen(ctx context.Context, w io.Writer, roomId string, flusher http.Flusher) error
}

type AdminWaitingRoomUseCase struct {
	rooms r.Rooms
	hub   IHub
}

var _ UseCaseInterface = (*AdminWaitingRoomUseCase)(nil)

func New(rooms r.Rooms, hub IHub) *AdminWaitingRoomUseCase {
	return &AdminWaitingRoomUseCase{rooms, hub}
}

func (uc *AdminWaitingRoomUseCase) ConnectAndListen(ctx context.Context, w io.Writer, roomId string, flusher http.Flusher) error {
	logger, _ := Logging.RetrieveLogger(ctx)

	if _, found := uc.rooms[r.Id(roomId)]; !found {
		return fmt.Errorf("room does not exist")
	}

	ch := uc.hub.Subscribe(roomId, "admin")
	defer func() {
		logger.Error("AdminWaitingRoomUseCase.ConnectAndListen: Unsubscribing...")
		if err := uc.hub.Unsubscribe(roomId, "admin"); err != nil {
			logger.Error("AdminWaitingRoomUseCase.ConnectAndListen:Error Unsubscribing", slog.Any("Error", err))
		}
	}()

	connectedMessage := fmt.Sprintf(`{"MessageType":"Connection","Message":"Connected to the room %s"}`, roomId)
	_, err := fmt.Fprintf(w, "data: %s\n\n", connectedMessage)
	if err != nil {
		logger.ErrorContext(ctx, "Could not send connected message")
		return fmt.Errorf("could not write message, %s", connectedMessage)
	}
	flusher.Flush()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-ch:
			if !ok {
				logger.ErrorContext(ctx, fmt.Sprintf("AdminWaitingRoomUseCase.ConnectAndListen: Could not receive the message in the room, %s", roomId))
				return nil
			}

			msg := string(message)
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				logger.ErrorContext(ctx, "Could not send message")
			}
			flusher.Flush()
		case <-ticker.C:
			_, err := fmt.Fprintf(w, "data: %s\n\n", `{"MessageType":"ping"}`)
			if err != nil {
				logger.ErrorContext(ctx, "Could not send ping")
			}
			flusher.Flush()
		case <-ctx.Done():
			return errors.New("context done?")
		}
	}
}
