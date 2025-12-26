package Messages

import (
	"ChoHanJi/domain/Item"
	"ChoHanJi/domain/Map"
	"ChoHanJi/domain/Player"
)

type ConnectedMessage struct {
	MapHeight int              `json:"MapHeight"`
	MapWidth  int              `json:"MapWidth"`
	Tiles     []*Map.Tile      `json:"Tiles"`
	Players   []*Player.Struct `json:"Players"`
	Items     []*Item.Struct   `json:"Items"`
}
