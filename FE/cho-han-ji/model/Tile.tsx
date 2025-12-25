import Item from "./Item";
import Player from "./Player";

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
  private _players: Record<string, Player>;
  private _items: Record<string, Item>;
  private _flag: Flag;
  private _team: Teams;

  constructor(x: number, y: number)
  constructor(x: number, y: number, flagInput: Flag)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>, items?: Record<string, Item>)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>, items?: Record<string, Item>, teams?: Teams) {
    if (x < 0 || y < 0) {
      throw new Error("Tile out of range")
    }

    this._flag = flagInput ?? Flag.EMPTY;
    this._players = players ?? {};
    this._items = items ?? {};
    this._team = teams ?? Teams.UNOCCUPIED;
  }

  get Flag(): Flag {
    return this._flag;
  }

  get Players(): Record<string, Player> {
    return this._players;
  }

  get Items(): Record<string, Item> {
    return this._items;
  }

  get Teams(): Teams {
    return this._team;
  }
}
