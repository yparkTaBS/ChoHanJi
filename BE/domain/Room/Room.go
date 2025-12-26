package Room

import (
	"ChoHanJi/domain/IdGenerator"
	"ChoHanJi/domain/Item"
	m "ChoHanJi/domain/Map"
	"ChoHanJi/domain/Player"
)

type Room struct {
	Map     *m.Map
	Players map[Player.Id]*Player.Struct
	Items   map[Item.Id]*Item.Struct
}

type (
	Id    string
	Rooms map[Id]*Room
)

func New() Rooms {
	return make(map[Id]*Room)
}

func CreateRoom(rooms Rooms, fieldMap *m.Map) (Id, error) {
	var id Id
	for {
		strId, err := IdGenerator.NewId()
		if err != nil {
			return "", err
		}

		id = Id(strId)
		_, found := rooms[id]
		if found {
			continue
		}

		room := new(Room)
		room.Map = fieldMap
		room.Players = make(map[Player.Id]*Player.Struct)
		room.Items = make(map[Item.Id]*Item.Struct)

		rooms[id] = room

		break
	}

	return id, nil
}
