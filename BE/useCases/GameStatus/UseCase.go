package GameStatus

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"ChoHanJi/infrastructure/Logging"
	"ChoHanJi/useCases/GameStatus/Messages"
	"context"
	"encoding/json"
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
		return fmt.Errorf("GameStatusUseCase.ConnectAndListen: %s %w", "room", ErrNotFound)
	}

	_, found = room.Players[Player.Id(playerId)]
	if playerId != "admin" && !found {
		return fmt.Errorf("GameStatusUseCase.ConnectAndListen: %s %w", "player", ErrNotFound)
	}

	ch := g.roomHub.Subscribe(roomId, playerId)
	defer func() {
		if err := g.roomHub.Unsubscribe(roomId, playerId); err != nil {
			logger.Error("GameStatusUseCase.ConnectAndListen:Error Unsubscribing", slog.Any("Error", err))
		}
	}()

	msgBody, err := g.getConnectedMessage(room)
	if err != nil {
		return err
	}

	connectedMessage := fmt.Sprintf(`{"MessageType":"Connection","Message":%s}`, string(msgBody))
	_, err = fmt.Fprintf(w, "data: %s\n\n", connectedMessage)
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
				logger.ErrorContext(ctx, fmt.Sprintf("GameStatusUseCase.ConnectAndListen: Could not receive the message in the room, %s", roomId))
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

func (g *UseCase) getConnectedMessage(room *Room.Room) ([]byte, error) {
	height, err := room.Map.GetMapHeight()
	if err != nil {
		return nil, err
	}

	width, err := room.Map.GetMapWidth()
	if err != nil {
		return nil, err
	}

	var players []*Player.Struct
	for _, val := range room.Players {
		players = append(players, val)
	}

	var items []*Item.Struct
	for _, val := range room.Items {
		items = append(items, val)
	}

	message := Messages.ConnectedMessage{
		MapHeight: height,
		MapWidth:  width,
		Tiles:     room.Map.GetRelevantTiles(),
		Players:   players,
		Items:     items,
	}

	return json.Marshal(message)
}
