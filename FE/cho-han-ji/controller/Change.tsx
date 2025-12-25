import Item from "@/model/Item";
import Player from "@/model/Player";

export type ObjectType = "Player" | "Item";

export default class Change {
  public X: number;
  public Y: number;
  public PrevX: number;
  public PrevY: number;
  public ObjectType: ObjectType;
  public Object: Player | Item;

  constructor(x: number, y: number, prevX: number, prevY: number, objectType: ObjectType, object: Player | Item) {
    this.X = x;
    this.Y = y;
    this.PrevX = prevX;
    this.PrevY = prevY;
    this.ObjectType = objectType;
    this.Object = object
  }

  public static fromJSON(data: any): Change {
    let object: Player | Item;
    if (data.ObjectType === "Player") {
      object = Player.fromJSON(data.Object);
    } else {
      object = Item.fromJSON(data.Object);
    }

    return new Change(Number(data.X), Number(data.Y), Number(data.PrevX), Number(data.PrevY), data.ObjectType as ObjectType, object);
  }
}
