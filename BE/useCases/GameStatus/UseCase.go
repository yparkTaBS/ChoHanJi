package GameStatus

import (
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"ChoHanJi/infrastructure/Logging"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var ErrNotFound error = errors.New("not found")

type Interface interface {
	ConnectAndListen(ctx context.Context, w io.Writer, roomId string, playerId string, flusher http.Flusher) error
}

type IHub interface {
	Subscribe(roomId, subscriberId string) <-chan []byte
	Unsubscribe(roomId, subscriberId string) error
}

type UseCase struct {
	rooms   Room.Rooms
	roomHub IHub
}

var _ Interface = (*UseCase)(nil)

func New(rooms Room.Rooms, roomHub IHub) *UseCase {
	return &UseCase{rooms, roomHub}
}

// ConnectAndListen implements GameStatusInterface.
func (g *UseCase) ConnectAndListen(ctx context.Context, w io.Writer, roomId string, playerId string, flusher http.Flusher) error {
	logger, _ := Logging.RetrieveLogger(ctx)

	room, found := g.rooms[Room.Id(roomId)]
	if !found {
		return fmt.Errorf("PlayerWaitingRoomUseCase.ConnectAndListen: %s %w", "room", ErrNotFound)
	}

	_, found = room.Players[Player.Id(playerId)]
	if !found {
		return fmt.Errorf("PlayerWaitingRoomUseCase.ConnectAndListen: %s %w", "player", ErrNotFound)
	}

	ch := g.roomHub.Subscribe(roomId, playerId)
	defer func() {
		logger.Error("PlayerWaitingRoomUseCase.ConnectAndListen: Unsubscribing...")
		if err := g.roomHub.Unsubscribe(roomId, playerId); err != nil {
			logger.Error("PlayerWaitingRoomUseCase.ConnectAndListen:Error Unsubscribing", slog.Any("Error", err))
		}
	}()

	connectedMessage := fmt.Sprintf(`{"MessageType":"Connection","Message":"Connected to the room %s"}`, roomId)
	_, err := fmt.Fprintf(w, "data: %s\n\n", connectedMessage)
	if err != nil {
		logger.ErrorContext(ctx, "Could not send connected message")
		return fmt.Errorf("could not write message, %s", connectedMessage)
	}
	flusher.Flush()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-ch:
			if !ok {
				logger.ErrorContext(ctx, fmt.Sprintf("PlayerWaitingRoomUseCase.ConnectAndListen: Could not receive the message in the room, %s", roomId))
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
			return nil
		}
	}
}
