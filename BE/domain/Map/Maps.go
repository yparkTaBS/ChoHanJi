package Map

import (
	"ChoHanJi/domain/IdGenerator"
	"errors"
)

var ErrMapNotFound error = errors.New("map not found")

type (
	Id   string
	Maps map[Id]*Map
)

func New() Maps {
	return make(map[Id]*Map)
}

func CreateMap(maps Maps, width, height int, items string) (Id, error) {
	var id Id
	for {
		strId, err := IdGenerator.NewId()
		if err != nil {
			return "", err
		}

		id = Id(strId)
		_, found := maps[id]
		if found {
			continue
		}

		fieldMap, err := NewMap(width, height, items)
		if err != nil {
			return "", err
		}
		maps[id] = fieldMap
		break
	}

	return id, nil
}

func GetMap(maps Maps, id Id) (*Map, error) {
	fieldMap, found := maps[id]
	if !found {
		return nil, ErrMapNotFound
	}
	return fieldMap, nil
}
