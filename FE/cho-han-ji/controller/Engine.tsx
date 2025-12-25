import Tile, { Flag, Teams } from "@/model/Tile";
import Change from "../model/Change";
import Player from "@/model/Player";
import Item from "@/model/Item";

export default class Engine {
  private tiles: Tile[][];
  private players: Record<string, Player>;
  private items: Record<string, Item>;

  constructor(width: number, height: number) {
    this.tiles = []
    for (let y = 0; y < height; y++) {
      let row: Tile[] = [];
      for (let x = 0; x < width; x++) {
        row.push(new Tile(x, y, Flag.EMPTY, {}, {}))
      }
      this.tiles.push(row);
    }

    this.players = {};
    this.items = {};
  }

  public Initialize(tiles: Tile[], players: Player[], items: Item[]) {
    for (const t of tiles) {
      this.tiles[t.Y][t.X].Flag = t.Flag;
      this.tiles[t.Y][t.X].Teams = t.Teams;
    }

    for (const p of players) {
      this.tiles[p.Y][p.X].Players[p.Id] = p;
      this.players[p.Id] = p;
    }

    for (const i of items) {
      this.tiles[i.Y][i.X].Items[i.Id] = i;
      this.items[i.Id] = i;
    }
  }

  public Update(changes: Change[]) {
    for (const change of changes) {
      const prevTile = this.tiles[change.PrevY][change.PrevX];
      if (change.X === -1 && change.Y === -1) {
        delete prevTile.Items[change.Id];
        continue;
      }

      const tile = this.tiles[change.Y][change.X]
      if (change.ObjectType === "Player") {
        delete prevTile.Players[change.Id];
        tile.Players[change.Id] = this.players[change.Id];
      } else if (change.ObjectType === "Item") {
        delete prevTile.Items[change.Id];
        tile.Items[change.Id] = this.items[change.Id];
      }
    }
  }

  public RenderAll(): [string, string, Flag, Teams][][] {
    return this.RenderParts(0, 0, this.tiles.length, this.tiles[0].length, true)
  }

  public RenderParts(startX: number, startY: number, width: number, height: number, condensed: boolean): [string, string, Flag, Teams][][] {
    if (startX < 0 || startY < 0) {
      throw new Error("Index out of Range");
    }

    const grid = this.CreateGrid(width, height);

    const maxY = this.tiles.length;
    const maxX = this.tiles[0].length;

    for (let x = 0; x < grid.length; x++) {
      for (let y = 0; y < grid[0].length; y++) {
        const rX = x + startX
        const rY = y + startY

        if ((rX >= maxX || rY >= maxY) || (rX < 0 || rY < 0)) {
          grid[x][y] = ["", "", Flag.INACCESSIBLE, Teams.UNOCCUPIED];
          continue;
        }

        const players = Object.entries(this.tiles[rY][rX].Players).map(([_, p]) => {
          if (!condensed) {
            return p.Name
          }
          return p.Name[0]
        }).join(",");
        const counts = Object.values(this.tiles[rY][rX].Items).reduce<Record<string, number>>((acc, item) => {
          const name = condensed ? item.Name[0] : item.Name;
          acc[name] = (acc[name] ?? 0) + 1;
          return acc;
        }, {});
        const items = Object.entries(counts).map(([key, val]) => `${key}x${val}`).join("|")
        const flag = this.tiles[rY][rX].Flag;
        const teams = this.tiles[rY][rX].Teams;

        grid[x][y] = [players, items, flag, teams]
      }
    }

    return grid;
  }

  private CreateGrid(width: number, height: number): [string, string, Flag, Teams][][] {
    let grid: [string, string, Flag, Teams][][] = []
    for (let x = 0; x < width; x++) {
      var row: [string, string, Flag, Teams][] = [];
      for (let y = 0; y < height; y++) {
        var tile: [string, string, Flag, Teams] = ["", "", Flag.EMPTY, Teams.UNOCCUPIED]
        row.push(tile)
      }
      grid.push(row);
    }
    return grid;
  }
}
