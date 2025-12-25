export enum PlayerClass {
  Fighter = 0,
  Archer = 1,
  Rogue = 2
}

export default class Player {
  private _id: string;
  private _name: string;
  private _playerClass: PlayerClass;
  private _bag?: string;

  constructor(id: string, name: string, playerClass: PlayerClass) {
    this._id = id;
    this._name = name;
    this._playerClass = playerClass;
    this._bag = "";
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

  get Bag(): string | undefined {
    return this._bag;
  }
}
