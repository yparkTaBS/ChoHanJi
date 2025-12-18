package RoomFactory

import (
	m "ChoHanJi/domain/Map"
	"ChoHanJi/useCases/RoomFactory/ports"
	"errors"
	"fmt"
)

type RoomFactory struct {
	maps m.Maps
}

var _ ports.IRoomFactory = (*RoomFactory)(nil)

func NewRoomFactory(maps m.Maps) (*RoomFactory, error) {
	if maps == nil {
		return nil, errors.New("Room.NewRoomFactory: maps data is null")
	}
	return &RoomFactory{maps}, nil
}

func (f *RoomFactory) Create(width, height int, items string) (m.Id, error) {
	id, err := m.CreateMap(f.maps, width, height, items)
	if err != nil {
		return "", fmt.Errorf("RoomFactoryi.Create: Failed to create the room: %w", err)
	}
	return id, nil
}
