package Map

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Team"
	"ChoHanJi/domain/TileFlag"
	"errors"
	"math/rand"
)

type Map struct {
	tiles [][]*Tile
}

func NewMap(width, height int, items []*Item.Struct) (*Map, error) {
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
			fieldMap.tiles[x][y] = NewTile(x, y, 0)
		}
	}

	// Set the Flags
	fieldMap.tiles[1][1].Flag = TileFlag.TREASURE_CHEST
	fieldMap.tiles[1][1].Team = Team.Team1
	fieldMap.tiles[width-2][1].Flag = TileFlag.SPAWN
	fieldMap.tiles[width-2][1].Team = Team.Team1
	fieldMap.tiles[width-2][height-2].Flag = TileFlag.TREASURE_CHEST
	fieldMap.tiles[width-2][height-2].Team = Team.Team2
	fieldMap.tiles[1][height-2].Flag = TileFlag.SPAWN
	fieldMap.tiles[1][height-2].Team = Team.Team2

	empty := fieldMap.getEmptyTileCoords()
	if len(empty) == 0 && len(items) > 0 {
		return nil, errors.New("no empty tiles available for item placement")
	}

	for _, item := range items {
		xy := empty[rand.Intn(len(empty))]
		x, y := xy[0], xy[1]

		item.X = x
		item.Y = y
		fieldMap.tiles[x][y].AddItem(item)
	}

	return fieldMap, nil
}

func (m *Map) PlacePlayer(player *Player.Struct) error {
	if m == nil || !m.IsInitialized() {
		return errors.New("map not initialized")
	}
	if player == nil {
		return errors.New("player is nil")
	}

	spawnTile, err := m.findSpawnTile(Team.Enum(player.TeamNumber))
	if err != nil {
		return err
	}

	spawnTile.AddPlayer(player)
	player.X = spawnTile.X
	player.Y = spawnTile.Y
	return nil
}

func (m *Map) findSpawnTile(team Team.Enum) (*Tile, error) {
	for _, column := range m.tiles {
		for _, tile := range column {
			if tile.Flag == TileFlag.SPAWN && tile.Team == team {
				return tile, nil
			}
		}
	}

	return nil, errors.New("spawn point not found for team")
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

func (m *Map) GetSpecialTiles() []*Tile {
	var tiles []*Tile
	for row := range m.tiles {
		for column := range m.tiles[row] {
			tile := m.tiles[row][column]
			if tile.IsSpecial() {
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
}

func (m *Map) GetNonEmptyTiles() []*Tile {
	var tiles []*Tile
	for row := range m.tiles {
		for column := range m.tiles[row] {
			tile := m.tiles[row][column]
			if tile.IsNotEmpty() {
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
}

func (m *Map) GetSpawn(team Team.Enum) (int, int) {
	width, _ := m.GetMapWidth()
	height, _ := m.GetMapHeight()
	if team == Team.Team1 {
		return width - 2, 1
	} else {
		return 1, height - 2
	}
}

func (m *Map) GetTeamTreasureChestLocation(team Team.Enum) (int, int) {
	width, _ := m.GetMapWidth()
	height, _ := m.GetMapHeight()
	if team == Team.Team1 {
		return 1, 1
	} else {
		return width - 2, height - 2
	}
}

func (m *Map) DisperseItems(team Team.Enum) ([]*Item.Struct, error) {
	x, y := m.GetTeamTreasureChestLocation(team)
	tile, err := m.GetTile(x, y)
	if err != nil {
		return nil, err
	}

	items := tile.Items
	if len(items) == 0 {
		return nil, nil
	}

	empty := m.getEmptyTileCoords()
	if len(empty) == 0 {
		return nil, errors.New("no empty tiles available for dispersing items")
	}

	var itemsMoved []*Item.Struct
	for _, item := range items {
		xy := empty[rand.Intn(len(empty))]
		x, y := xy[0], xy[1]

		item.X = x
		item.Y = y
		m.tiles[x][y].AddItem(item)
		itemsMoved = append(itemsMoved, item)
	}

	tile.Items = tile.Items[:0]
	return itemsMoved, nil
}

func (m *Map) getEmptyTileCoords() [][2]int {
	var coords [][2]int
	for x := range m.tiles {
		for y := range m.tiles[x] {
			if m.tiles[x][y].Flag == TileFlag.EMPTY {
				coords = append(coords, [2]int{x, y})
			}
		}
	}
	return coords
}
