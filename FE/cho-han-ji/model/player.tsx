import Item from "./Item";
import { Teams } from "./Tile";

export enum PlayerClass {
  Fighter = 0,
  Archer = 1,
  Rogue = 2
}

export default class Player {
  public X: number;
  public Y: number;
  public Id: string;
  public Name: string;
  public Class: PlayerClass;
  public Team: Teams
  public _bag?: Item;

  constructor(x: number, y: number, id: string, name: string, playerClass: PlayerClass, team: Teams) {
    this.X = x;
    this.Y = y;
    this.Id = id;
    this.Name = name;
    this.Class = playerClass;
    this.Team = team;
  }
}
