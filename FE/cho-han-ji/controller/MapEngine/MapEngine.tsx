import Tile, { Flag, Teams } from "@/model/Tile";

export default class Engine {
  private tiles: Tile[][];

  constructor(width: number, height: number) {
    this.tiles = []
    for (let y = 0; y < height; y++) {
      let row: Tile[] = [];
      for (let x = 0; x < width; x++) {
        row.push(new Tile(x, y, Flag.EMPTY, [], []))
      }
      this.tiles.push(row);
    }
  }

  public Update() { }

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

        const players = this.tiles[rY][rX].Players.map(p => {
          if (!condensed) {
            return p.Name
          }
          return p.Name[0]
        }).join(",");
        const items = Object.entries(this.tiles[rY][rX].Items.reduce<Record<string, number>>((acc, s) => {
          acc[s] = (acc[s] ?? 0) + 1;
          return acc;
        }, {})).map(([key, val]) => {
          if (condensed) {
            key = key[0];
          }
          return `${key}x${val}`;
        }).join("|")
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
