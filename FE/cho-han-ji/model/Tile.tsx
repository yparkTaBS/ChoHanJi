import Item from "./Item";
import Player from "./Player";

export enum Flag {
  EMPTY = 0,
  SPAWN = 1,
  TREASURE_CHEST = 2,
  INACCESSIBLE = 3,
}

export enum Teams {
  Neutral = 0,
  TEAM1 = 1,
  TEAM2 = 2
}

export default class Tile {
  public X: number;
  public Y: number;
  public Players: Record<string, Player>;
  public Items: Record<string, Item>;
  public Flag: Flag;
  public Team: Teams;

  constructor(x: number, y: number)
  constructor(x: number, y: number, flagInput: Flag)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>, items?: Record<string, Item>)
  constructor(x: number, y: number, flagInput?: Flag, players?: Record<string, Player>, items?: Record<string, Item>, teams?: Teams) {
    if (x < 0 || y < 0) {
      throw new Error("Tile out of range")
    }

    this.X = x;
    this.Y = y;
    this.Flag = flagInput ?? Flag.EMPTY;
    this.Players = players ?? {};
    this.Items = items ?? {};
    this.Team = teams ?? Teams.Neutral;
  }
}
