package PlayerWaitingRoomUseCase

import (
	"ChoHanJi/domain/Map"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"ChoHanJi/infrastructure/Logging"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IPlayerHub interface {
	Subscribe(roomId string) (<-chan []byte, int)
	Unsubscribe(roomId string, index int) error
}

type IRoomHub interface {
	Publish(roomId, messageType, messageBody string) error
}

var (
	_ IPlayerHub = (*SSEHub.SSEHub)(nil)
	_ IRoomHub   = (*SSEHub.SSEHub)(nil)
)

var ErrNotFound error = errors.New("not found")

type PlayerWaitingRoomUseCase struct {
	rooms     Room.Rooms
	playerHub IPlayerHub
	roomHub   IRoomHub
}

type IPlayerWaitingRoomUseCase interface {
	ConnectAndListen(ctx context.Context, w io.Writer, roomId string, playerId string, flusher http.Flusher) error
}

var _ IPlayerWaitingRoomUseCase = (*PlayerWaitingRoomUseCase)(nil)

func New(rooms Room.Rooms, playerHub IPlayerHub, roomHub IRoomHub) (*PlayerWaitingRoomUseCase, error) {
	return &PlayerWaitingRoomUseCase{rooms, playerHub, roomHub}, nil
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

	Map.PlacePlayer(*room.Map, player)

	ch, index := p.playerHub.Subscribe(roomId)
	defer func() {
		_ = p.playerHub.Unsubscribe(roomId, index)
	}()

	connectedMessage := fmt.Sprintf(`{"MessageType":"Connection","Message":"Connected to the room %s"}`, roomId)
	_, err := fmt.Fprintf(w, "data: %s\n\n", connectedMessage)
	if err != nil {
		logger.ErrorContext(ctx, "Could not send connected message")
		return fmt.Errorf("could not write message, %s", connectedMessage)
	}
	flusher.Flush()

	if err := p.roomHub.Publish(roomId, "PlayerConnected", fmt.Sprintf(`{"id":"%s","name":"%s"}`, playerId, player.Name)); err != nil {
		logger.ErrorContext(ctx, "Could not announce player connected message")
		return fmt.Errorf("could not publish PlayerConnected message. PlayerId: %s", player.Id)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-ch:
			if !ok {
				logger.ErrorContext(ctx, fmt.Sprintf("Could not receive the message in the room, %s", roomId))
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
