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
      this.tiles[t.Y][t.X].Team = t.Team;
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
      const prevTile =
        change.PrevX >= 0 &&
        change.PrevY >= 0 &&
        change.PrevY < this.tiles.length &&
        change.PrevX < this.tiles[0].length
          ? this.tiles[change.PrevY][change.PrevX]
          : null;

      const nextTile =
        change.X >= 0 &&
        change.Y >= 0 &&
        change.Y < this.tiles.length &&
        change.X < this.tiles[0].length
          ? this.tiles[change.Y][change.X]
          : null;

      if (change.X === -1 && change.Y === -1) {
        if (prevTile) {
          if (change.ObjectType === "Player") {
            delete prevTile.Players[change.Id];
          } else if (change.ObjectType === "Item") {
            delete prevTile.Items[change.Id];
          }
        }
        continue;
      }

      if (!nextTile) continue;

      if (change.ObjectType === "Player") {
        if (prevTile) {
          delete prevTile.Players[change.Id];
        }
        nextTile.Players[change.Id] = this.players[change.Id];
      } else if (change.ObjectType === "Item") {
        if (prevTile) {
          delete prevTile.Items[change.Id];
        }
        nextTile.Items[change.Id] = this.items[change.Id];
      }
    }
  }

  public RenderAll(): [string, string, Flag, Teams][][] {
    return this.RenderParts(0, 0, this.tiles.length, this.tiles[0].length, true)
  }

  public RenderParts(startX: number, startY: number, height: number, width: number, condensed: boolean): [string, string, Flag, Teams][][] {
    const grid = this.CreateGrid(width, height);

    const maxY = this.tiles.length;
    const maxX = this.tiles[0].length;

    for (let y = 0; y < grid.length; y++) {
      for (let x = 0; x < grid[0].length; x++) {
        const rX = x + startX
        const rY = y + startY

        if ((rX >= maxX || rY >= maxY) || (rX < 0 || rY < 0)) {
          grid[y][x] = ["", "", Flag.INACCESSIBLE, Teams.Neutral];
          continue;
        }

        const players = Object.entries(this.tiles[rY][rX].Players).map(entry => { return entry[1] })
        const playerNames = players.map(p => {
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
        var teams = this.tiles[rY][rX].Team;

        if (teams === Teams.Neutral) {
          if (players.length != 0) {
            teams = players[0].Team
          }
        }

        grid[y][x] = [playerNames, items, flag, teams]
      }
    }

    return grid;
  }

  public GetChestItemsByTeam(): Record<Teams, Record<string, number>> {
    const chestItems: Record<Teams, Record<string, number>> = {
      [Teams.Neutral]: {},
      [Teams.TEAM1]: {},
      [Teams.TEAM2]: {},
    };

    for (const row of this.tiles) {
      for (const tile of row) {
        if (tile.Flag !== Flag.TREASURE_CHEST) continue;
        const bucket = chestItems[tile.Team] ?? chestItems[Teams.Neutral];

        for (const item of Object.values(tile.Items)) {
          bucket[item.Name] = (bucket[item.Name] ?? 0) + 1;
        }
      }
    }

    return chestItems;
  }

  private CreateGrid(height: number, width: number): [string, string, Flag, Teams][][] {
    let grid: [string, string, Flag, Teams][][] = []
    for (let y = 0; y < height; y++) {
      var row: [string, string, Flag, Teams][] = [];
      for (let x = 0; x < width; x++) {
        var tile: [string, string, Flag, Teams] = ["", "", Flag.EMPTY, Teams.Neutral]
        row.push(tile)
      }
      grid.push(row);
    }
    return grid;
  }
}
