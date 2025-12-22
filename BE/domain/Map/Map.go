package Map

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
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

	items := strings.SplitSeq(itemList, ",")
	for itemName := range items {
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

func PlacePlayer(fieldMap Map, player *Player.Player) {
	if player.TeamNumber == 1 {
		fieldMap.tiles[0][0].AddPlayer(player)
	} else {
		x := len(fieldMap.tiles) - 1
		y := len(fieldMap.tiles[0]) - 1
		fieldMap.tiles[x][y].AddPlayer(player)
	}
}
