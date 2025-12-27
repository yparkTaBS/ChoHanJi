package StartGameUseCase

import (
	"ChoHanJi/domain/Action"
	"ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"encoding/json"
	"errors"
)

type IActionList interface {
	StartGame(roomId Room.Id)
}

type IHub interface {
	PublishToAll(roomId, messageType, messageBody string) error
}

var (
	_ IActionList = (*Action.List)(nil)
	_ IHub        = (*SSEHub.Struct)(nil)
)

type Interface interface {
	Announce(roomId string) error
}

type Struct struct {
	rooms Room.Rooms
	list  IActionList
	hub   IHub
}

var _ Interface = (*Struct)(nil)

func New(rooms Room.Rooms, list IActionList, hub IHub) *Struct {
	return &Struct{rooms, list, hub}
}

// Announce implements IStartGameUseCase.
func (s *Struct) Announce(roomId string) error {
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

	s.list.StartGame(Room.Id(roomId))

	return nil
}

type MessageBody struct {
	RoomId    string `json:"RoomId"`
	MapHeight int    `json:"MapHeight"`
	MapWeidth int    `json:"MapWidth"`
}
