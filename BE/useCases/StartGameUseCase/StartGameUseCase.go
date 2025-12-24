package StartGameUseCase

import (
	"ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"encoding/json"
	"errors"
)

type UseCaseInterface interface {
	Announce(roomId string) error
}

type StartGameUseCase struct {
	rooms Room.Rooms
	hub   IHub
}

var _ UseCaseInterface = (*StartGameUseCase)(nil)

type IHub interface {
	PublishToAll(roomId, messageType, messageBody string) error
}

var _ IHub = (*SSEHub.SSEHub)(nil)

func New(rooms Room.Rooms, hub IHub) *StartGameUseCase {
	return &StartGameUseCase{rooms, hub}
}

// Announce implements IStartGameUseCase.
func (s *StartGameUseCase) Announce(roomId string) error {
	room, found := s.rooms[Room.Id(roomId)]
	if !found {
		return errors.New("game not found")
	}

	height, err := room.Map.GetMapHeight()
	if err != nil {
		return err
	}

	width, err := room.Map.GetMapWidth()
	if err != nil {
		return err
	}

	messageBody := MessageBody{roomId, height, width}

	msg, err := json.Marshal(messageBody)
	if err != nil {
		return errors.New("failed to send the game start message")
	}

	if err := s.hub.PublishToAll(roomId, "GameStart", string(msg)); err != nil {
		return err
	}

	return nil
}

type MessageBody struct {
	RoomId    string `json:"RoomId"`
	MapHeight int    `json:"MapHeight"`
	MapWeidth int    `json:"MapWidth"`
}
