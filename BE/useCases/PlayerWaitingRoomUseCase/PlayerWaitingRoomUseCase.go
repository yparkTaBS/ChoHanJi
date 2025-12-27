package PlayerWaitingRoomUseCase

import (
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
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
	Publish(roomId, adminId, messageType, messageBody string) error
}

var _ IHub = (*SSEHub.SSEHub)(nil)

var ErrNotFound error = errors.New("not found")

type PlayerWaitingRoomUseCase struct {
	rooms   Room.Rooms
	roomHub IHub
}

type UseCaseInterface interface {
	ConnectAndListen(ctx context.Context, w io.Writer, roomId string, playerId string, flusher http.Flusher) error
}

var _ UseCaseInterface = (*PlayerWaitingRoomUseCase)(nil)

func New(rooms Room.Rooms, roomHub IHub) (*PlayerWaitingRoomUseCase, error) {
	return &PlayerWaitingRoomUseCase{rooms, roomHub}, nil
}

func (p *PlayerWaitingRoomUseCase) ConnectAndListen(ctx context.Context, w io.Writer, roomId string, playerId string, flusher http.Flusher) error {
	logger, _ := Logging.RetrieveLogger(ctx)

	room, found := p.rooms[Room.Id(roomId)]
	if !found {
		return fmt.Errorf("PlayerWaitingRoomUseCase.ConnectAndListen: %s %w", "room", ErrNotFound)
	}

	player, found := room.Players[Player.Id(playerId)]
	if !found {
		return fmt.Errorf("PlayerWaitingRoomUseCase.ConnectAndListen: %s %w", "player", ErrNotFound)
	}

	room.Map.PlacePlayer(player)

	ch := p.roomHub.Subscribe(roomId, playerId)
	defer func() {
		logger.Error("PlayerWaitingRoomUseCase.ConnectAndListen: Unsubscribing...")
		if err := p.roomHub.Unsubscribe(roomId, playerId); err != nil {
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

	if err := p.roomHub.Publish(roomId, "admin", "PlayerConnected", fmt.Sprintf(`{"id":"%s","name":"%s","team":%d}`, playerId, player.Name, player.TeamNumber)); err != nil {
		logger.ErrorContext(ctx, "PlayerWaitingRoomUseCase.ConnectAndListen: Could not announce player connected message", slog.Any("Error", err))
		return fmt.Errorf("could not publish PlayerConnected message. PlayerId: %s", player.Id)
	}

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
