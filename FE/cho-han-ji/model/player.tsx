import Item from "./Item";

export enum PlayerClass {
  Fighter = 0,
  Archer = 1,
  Rogue = 2
}

export default class Player {
  private _x: number;
  private _y: number;
  private _id: string;
  private _name: string;
  private _playerClass: PlayerClass;
  private _bag?: Item;

  constructor(x: number, y: number, id: string, name: string, playerClass: PlayerClass) {
    this._x = x;
    this._y = y;
    this._id = id;
    this._name = name;
    this._playerClass = playerClass;
  }

  get X(): number {
    return this._x;
  }

  set X(x) {
    this._x = x;
  }

  get Y(): number {
    return this._y;
  }

  set Y(y) {
    this._y = y;
  }

  get Id(): string {
    return this._id;
  }

  get Name(): string {
    return this._name;
  }

  get Class(): PlayerClass {
    return this._playerClass;
  }

  get Bag(): Item | undefined {
    return this._bag;
  }

  set Bag(item: Item | undefined) {
    this._bag = item;
  }

  public static fromJSON(data: any): Player {
    const item = data.Bag ? Item.fromJSON(data.Bag) : undefined;

    const player = new Player(
      Number(data.X),
      Number(data.Y),
      String(data.Id),
      String(data.Name),
      Number(data.Class) as PlayerClass,
    );

    if (item) {
      player.Bag = item;
    }

    return player;
  }
}
