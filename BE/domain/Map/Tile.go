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

func (t *Tile) IsRelevant() bool {
	return t.Flag != TileFlag.EMPTY
}
