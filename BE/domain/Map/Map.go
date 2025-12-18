package Map

import (
	"ChoHanJi/domain/Item"
	"math/rand"
	"strings"
)

type Map struct {
	tiles [][]*Tile
}

func NewMap(width, height int, itemList string) (*Map, error) {
	fieldMap := new(Map)
	fieldMap.tiles = make([][]*Tile, width)
	for i := range fieldMap.tiles {
		fieldMap.tiles[i] = make([]*Tile, height)
	}

	for x := range width {
		for y := range height {
			fieldMap.tiles[x][y] = NewTile()
		}
	}

	items := strings.Split(itemList, ",")
	for _, itemName := range items {
		x := rand.Intn(width)
		y := rand.Intn(height)
		item, err := Item.New(itemName)
		if err != nil {
			return nil, err
		}
		fieldMap.tiles[x][y].AddItem(item)
	}

	return fieldMap, nil
}
