import Player from "./player";

export enum Flag {
  EMPTY = 0,
  SPAWN = 1,
  TREASURE_CHEST = 2,
  INACCESSIBLE = 3,
}

export enum Teams {
  UNOCCUPIED = 0,
  TEAM1 = 1,
  TEAM2 = 2
}

export default class Tile {
  private x: number;
  private y: number;
  private _players: Player[];
  private _items: string[];
  private _flag: Flag;
  private _team: Teams;

  constructor(x: number, y: number)
  constructor(x: number, y: number, flagInput: Flag)
  constructor(x: number, y: number, flagInput?: Flag, players?: Player[])
  constructor(x: number, y: number, flagInput?: Flag, players?: Player[], items?: string[])
  constructor(x: number, y: number, flagInput?: Flag, players?: Player[], items?: string[], teams?: Teams) {
    if (x < 0 || y < 0) {
      throw new Error("Tile out of range")
    }

    this.x = x;
    this.y = y;

    this._flag = flagInput ?? Flag.EMPTY;
    this._players = players ?? [];
    this._items = items ?? [];
    this._team = teams ?? Teams.UNOCCUPIED;
  }

  get Flag(): Flag {
    return this._flag;
  }

  get Players(): Player[] {
    return this._players;
  }

  get Items(): string[] {
    return this._items;
  }

  get Teams(): Teams {
    return this._team;
  }
}
