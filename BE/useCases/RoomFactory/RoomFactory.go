package RoomFactory

import (
	"ChoHanJi/domain/Item"
	m "ChoHanJi/domain/Map"
	r "ChoHanJi/domain/Room"
	"ChoHanJi/useCases/RoomFactory/ports"
	"errors"
	"fmt"
)

type RoomFactory struct {
	rooms r.Rooms
}

var _ ports.UseCaseInterface = (*RoomFactory)(nil)

func NewRoomFactory(rooms r.Rooms) (*RoomFactory, error) {
	if rooms == nil {
		return nil, errors.New("Room.NewRoomFactory: rooms data is null")
	}
	return &RoomFactory{rooms}, nil
}

func (f *RoomFactory) Create(width, height int, itemNames string) (r.Id, error) {
	items := Item.DecodeItems(itemNames)
	fieldMap, err := m.NewMap(width, height, items)
	if err != nil {
		return "", fmt.Errorf("RoomFactory.Create: Failed to create the map: %w", err)
	}

	id, err := r.CreateRoom(f.rooms, fieldMap)
	if err != nil {
		return "", fmt.Errorf("RoomFactory.Create: Failed to create the room %w", err)
	}

	for _, val := range items {
		f.rooms[id].Items[val.Id] = val
	}

	return id, nil
}
