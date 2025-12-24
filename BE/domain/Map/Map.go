package Map

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
	"errors"
	"math/rand"
	"strings"
)

type Map struct {
	tiles [][]*Tile
}

func NewMap(width, height int, itemList string) (*Map, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be > 0")
	}

	fieldMap := &Map{
		tiles: make([][]*Tile, width),
	}

	// Allocate 2D slice
	for x := range width {
		fieldMap.tiles[x] = make([]*Tile, height)
		for y := range height {
			fieldMap.tiles[x][y] = NewTile()
		}
	}

	// Add items if provided
	itemList = strings.TrimSpace(itemList)
	if itemList != "" {
		items := strings.SplitSeq(itemList, ",")
		for itemName := range items {
			itemName = strings.TrimSpace(itemName)
			if itemName == "" {
				continue
			}

			x := rand.Intn(width)
			y := rand.Intn(height)

			item, err := Item.New(itemName)
			if err != nil {
				return nil, err
			}

			fieldMap.tiles[x][y].AddItem(item)
		}
	}

	return fieldMap, nil
}

func (m *Map) PlacePlayer(player *Player.Player) error {
	if m == nil || !m.IsInitialized() {
		return errors.New("map not initialized")
	}
	if player == nil {
		return errors.New("player is nil")
	}

	if player.TeamNumber == 1 {
		m.tiles[0][0].AddPlayer(player)
		return nil
	}

	x := len(m.tiles) - 1
	y := len(m.tiles[0]) - 1
	m.tiles[x][y].AddPlayer(player)
	return nil
}

func (m *Map) GetMapWidth() (int, error) {
	if m == nil || !m.IsInitialized() {
		return 0, errors.New("map not initialized")
	}
	return len(m.tiles), nil
}

func (m *Map) GetMapHeight() (int, error) {
	if m == nil || !m.IsInitialized() {
		return 0, errors.New("map not initialized")
	}
	return len(m.tiles[0]), nil
}

func (m *Map) IsInitialized() bool {
	if m == nil {
		return false
	}
	if len(m.tiles) == 0 {
		return false
	}
	if len(m.tiles[0]) == 0 {
		return false
	}
	return true
}

func (m *Map) GetTile(x, y int) (*Tile, error) {
	if !m.IsInitialized() {
		return nil, errors.New("map not initialized")
	}
	if x < 0 || x >= len(m.tiles) {
		return nil, errors.New("x out of bounds")
	}
	if y < 0 || y >= len(m.tiles[0]) {
		return nil, errors.New("y out of bounds")
	}
	return m.tiles[x][y], nil
}
