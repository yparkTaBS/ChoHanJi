import Item from "../Item"
import Player from "../Player"
import Tile from "../Tile"

export class GameConnected {
  public MapHeight: number;
  public MapWidth: number;
  public Tiles: Tile[];
  public Players: Player[];
  public Items: Item[];

  constructor(MapHeight: number, MapWidth: number, tiles: Tile[], players: Player[], items: Item[]) {
    this.MapHeight = MapHeight;
    this.MapWidth = MapWidth;
    this.Tiles = tiles;
    this.Players = players;
    this.Items = items;
  }
}
