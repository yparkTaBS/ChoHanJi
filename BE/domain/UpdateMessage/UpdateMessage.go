package UpdateMessage

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
)

type Struct struct {
	PlayerChanges map[Player.Id]*PlayerChange `json:"PlayerChanges"`
	ItemChanges   map[Item.Id]*ItemChange     `json:"ItemChanges"`
}

func New() *Struct {
	return &Struct{
		make(map[Player.Id]*PlayerChange),
		make(map[Item.Id]*ItemChange),
	}
}

type PlayerChange struct {
	X      int
	Y      int
	PrevX  int
	PrevY  int
	Id     Player.Id
	ItemId *Item.Id
}

type ItemChange struct {
	X      int
	Y      int
	PrevX  int
	PrevY  int
	ItemId Item.Id
}

func (s *Struct) UpsertPlayer(id Player.Id, X, Y, PrevX, PrevY int, itemId *Item.Id) {
	if player, found := s.PlayerChanges[id]; !found {
		s.PlayerChanges[id] = &PlayerChange{X, Y, PrevX, PrevY, id, itemId}
	} else {
		player.X = X
		player.Y = Y
		player.ItemId = itemId
	}
}

func (s *Struct) UpsertItem(id Item.Id, X, Y, PrevX, PrevY int) {
	if item, found := s.ItemChanges[id]; !found {
		s.ItemChanges[id] = &ItemChange{X, Y, PrevX, PrevY, id}
	} else {
		item.X = X
		item.Y = Y
	}
}
