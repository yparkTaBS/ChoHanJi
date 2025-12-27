package Changes

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Player"
)

type Struct struct {
	PlayerChanges []PlayerChange `json:"PlayerChanges"`
	ItemChanges   []ItemChange   `json:"ItemChanges"`
}

type PlayerChange struct {
	X      int
	Y      int
	PrevX  int
	PrevY  int
	Id     Player.Id
	ItemId Item.Id
}

type ItemChange struct {
	X      int
	Y      int
	PrevX  int
	PrevY  int
	ItemId Item.Id
}
