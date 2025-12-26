package Player

import (
	c "ChoHanJi/domain/Class"
	"ChoHanJi/domain/IdGenerator"
	"ChoHanJi/domain/Item"
	"errors"
	"strings"
)

type Id string

type Struct struct {
	X          int `json:"X"`
	Y          int `json:"Y"`
	Id         Id
	IdStr      string `json:"Id"`
	Name       string `json:"Name"`
	Class      c.Struct
	ClassName  string       `json:"Class"`
	Bag        *Item.Struct `json:"Item,omitempty"`
	TeamNumber int          `json:"Team"`
}

func New(players map[Id]*Struct, name, class string, team int) (*Struct, error) {
	var player *Struct
	for {
		strId, err := IdGenerator.NewId()
		if err != nil {
			return nil, err
		}

		id := Id(strId)
		_, found := players[id]
		if found {
			continue
		}

		playerClass, err := getClass(strings.ToUpper(class))
		if err != nil {
			return nil, err
		}

		player = &Struct{
			Id:         id,
			IdStr:      string(id),
			Name:       name,
			Class:      playerClass,
			ClassName:  class,
			TeamNumber: team,
		}

		break
	}

	return player, nil
}

func getClass(class string) (c.Struct, error) {
	switch class {
	case "FIGHTER":
		return c.Fighter, nil
	case "RANGER":
		return c.Ranger, nil
	case "THIEF":
		return c.Rogue, nil
	default:
		return c.Struct{}, errors.New("class not defined")
	}
}
