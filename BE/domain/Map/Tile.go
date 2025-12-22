package Map

import (
	i "ChoHanJi/domain/Item"
	p "ChoHanJi/domain/Player"
)

type Tile struct {
	Player []*p.Player
	Items  []i.Item
}

func NewTile() *Tile {
	var items []i.Item
	var players []*p.Player
	return &Tile{players, items}
}

func (t *Tile) AddItem(item i.Item) {
	t.Items = append(t.Items, item)
}

func (t *Tile) AddPlayer(player *p.Player) {
	t.Player = append(t.Player, player)
}
