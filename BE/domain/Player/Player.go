package Player

import (
	c "ChoHanJi/domain/Class"
	"ChoHanJi/domain/IdGenerator"
	i "ChoHanJi/domain/Item"
	"errors"
	"strings"
)

type Id string

type Player struct {
	Id         Id
	Name       string
	Class      c.Class
	Bag        i.Item
	TeamNumber int
}

func New(players map[Id]*Player, name, class string, team int) (*Player, error) {
	var player *Player
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

		player = &Player{
			Id:         id,
			Name:       name,
			Class:      playerClass,
			TeamNumber: team,
		}

		break
	}

	return player, nil
}

func getClass(class string) (c.Class, error) {
	switch class {
	case "FIGHTER":
		return c.Fighter, nil
	case "RANGER":
		return c.Ranger, nil
	case "THIEF":
		return c.Rogue, nil
	default:
		return c.Class{}, errors.New("class not defined")
	}
}
