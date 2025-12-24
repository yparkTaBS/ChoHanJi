import Tile, { Flag } from "@/model/Tile";

export default class Engine {
  private tiles: Tile[][];

  constructor(width: number, height: number) {
    this.tiles = []
    for (let y = 0; y < height; y++) {
      let row: Tile[] = []
      for (let x = 0; x < width; x++) {
        row.push(new Tile(x, y, Flag.EMPTY, string[]))
      }

    }

    return this;
  }

  public Update() { }

  public Render()
}
