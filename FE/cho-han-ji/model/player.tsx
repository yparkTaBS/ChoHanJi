import Item from "./Item";

export enum PlayerClass {
  Fighter = 0,
  Archer = 1,
  Rogue = 2
}

export default class Player {
  private _id: string;
  private _name: string;
  private _playerClass: PlayerClass;
  private _bag?: Item;

  constructor(id: string, name: string, playerClass: PlayerClass) {
    this._id = id;
    this._name = name;
    this._playerClass = playerClass;
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
