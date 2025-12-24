export enum Flag {
  EMPTY = 0,
  SPAWN = 1,
  TREASURE_CHEST = 2,
}

export default class Tile {
  private x: number;
  private y: number;
  private flag: Flag;
  private contains: string[];

  constructor(x: number, y: number)
  constructor(x: number, y: number, flagInput: Flag)
  constructor(x: number, y: number, flagInput?: Flag, contains?: string[]) {
    if (x < 0 || y < 0) {
      throw new Error("Tile out of range")
    }

    this.x = x;
    this.y = y;

    this.flag = flagInput ?? Flag.EMPTY;
    this.contains = contains ?? [];
  }
}
