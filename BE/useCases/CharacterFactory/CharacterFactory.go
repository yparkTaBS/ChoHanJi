package CharacterFactory

import (
	"ChoHanJi/domain/Player"
	r "ChoHanJi/domain/Room"
	"errors"
)

type ICharacterFactory interface {
	CreateCharacter(string, string, string, int) (string, error)
}

type CharacterFactory struct {
	rooms r.Rooms
}

var _ ICharacterFactory = (*CharacterFactory)(nil)

func New(rooms r.Rooms) *CharacterFactory {
	return &CharacterFactory{rooms}
}

// CreateCharacter implements ICharacterFactory.
func (c *CharacterFactory) CreateCharacter(roomId string, name, class string, teamNumber int) (string, error) {
	room, found := c.rooms[r.Id(roomId)]
	if !found {
		return "", errors.New("the game room does not exist")
	}

	player, err := Player.New(room.Players, name, class, teamNumber)
	if err != nil {
		return "", err
	}

	room.Players[player.Id] = player

	return string(player.Id), nil
}
