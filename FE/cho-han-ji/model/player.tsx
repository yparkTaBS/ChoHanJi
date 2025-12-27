import Item from "./Item";
import { Teams } from "./Tile";

export enum PlayerClass {
  Fighter = "fighter",
  Archer = "ranger",
  Thief = "thief"
}

class Class {
  private _power: number;
  private _defence: number;
  private _range: number;
  private _movementSpeed: number;
  private _initialHP: number;

  constructor(power: number, defence: number, range: number, movementSpeed: number, initialHP: number) {
    this._power = power;
    this._defence = defence;
    this._range = range;
    this._movementSpeed = movementSpeed;
    this._initialHP = initialHP;
  }

  static ConvertToClass(c: PlayerClass): Class {
    if (c === PlayerClass.Fighter) {
      return new Fighter()
    } else if (c === PlayerClass.Archer) {
      return new Archer()
    } else if (c === PlayerClass.Thief) {
      return new Thief()
    } else {
      throw new Error("Class Not defined")
    }
  }

  get Power(): number { return this._power; }
  get Defence(): number { return this._defence; }
  get Range(): number { return this._range; }
  get MovementSpeed(): number { return this._movementSpeed; }
  get InitialHP(): number { return this._initialHP; }
}

class Fighter extends Class {
  constructor() {
    super(2, 1, 1, 1, 2)
  }
}

class Archer extends Class {
  constructor() {
    super(2, 0, 2, 1, 2)
  }
}

class Thief extends Class {
  constructor() {
    super(2, 0, 1, 3, 2)
  }
}

export default class Player {
  public X: number;
  public Y: number;
  public Id: string;
  public Name: string;
  public Class: PlayerClass;
  private _classInfo: Class | undefined;
  public Team: Teams
  public _bag?: Item;

  constructor(x: number, y: number, id: string, name: string, playerClass: PlayerClass, team: Teams) {
    this.X = x;
    this.Y = y;
    this.Id = id;
    this.Name = name;
    this.Class = playerClass;
    this._classInfo = Class.ConvertToClass(this.Class)
    this.Team = team;
  }

  get ClassInfo(): Class {
    if (!this._classInfo) {
      this._classInfo = Class.ConvertToClass(this.Class)
    }
    return this._classInfo
  }
}

export class PlayerInstance extends Player {
  public CurrentHP: number;
  public CurrentX: number;
  public CurrentY: number;

  constructor(player: Player) {
    super(player.X, player.Y, player.Id, player.Name, player.Class, player.Team)
    this.CurrentHP = this.ClassInfo.InitialHP;
    this.CurrentX = this.X;
    this.CurrentY = this.Y;
  }
}

