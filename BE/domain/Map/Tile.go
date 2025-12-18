package Map

import (
	i "ChoHanJi/domain/Item"
	p "ChoHanJi/domain/Player"
)

type Tile struct {
	Player *p.Player
	Items  []i.Item
}

func NewTile() *Tile {
	var items []i.Item
	return &Tile{nil, items}
}

func (t *Tile) AddItem(item i.Item) {
	t.Items = append(t.Items, item)
}
