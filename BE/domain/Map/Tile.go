package Map

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Team"
	"ChoHanJi/domain/TileFlag"
)

type Tile struct {
	X      int
	Y      int
	Flag   TileFlag.TileFlagEnum `json:"Flag"`
	Player []*Player.Struct      `json:"Player,omitempty"`
	Items  []*Item.Struct        `json:"Items,omitempty"`
	Team   Team.Enum             `json:"Team"`
}

func NewTile(x, y int, team Team.Enum) *Tile {
	var items []*Item.Struct
	var players []*Player.Struct
	return &Tile{x, y, TileFlag.EMPTY, players, items, team}
}

func (t *Tile) AddItem(item *Item.Struct) {
	t.Items = append(t.Items, item)
}

func (t *Tile) AddPlayer(player *Player.Struct) {
	t.Player = append(t.Player, player)
}

func (t *Tile) IsSpecial() bool {
	return t.Flag != TileFlag.EMPTY
}

func (t *Tile) IsNotEmpty() bool {
	return len(t.Player) != 0 || len(t.Items) != 0
}

func (t *Tile) RemovePlayer(playerId Player.Id) {
	for i, pl := range t.Player {
		if pl.Id == playerId {
			last := len(t.Player) - 1
			t.Player[i] = t.Player[last]
			t.Player = t.Player[:last]
			return
		}
	}
}
