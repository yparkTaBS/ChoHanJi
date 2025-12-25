export type ObjectType = "Player" | "Item";

export default class Change {
  public X: number;
  public Y: number;
  public PrevX: number;
  public PrevY: number;
  public ObjectType: ObjectType;
  public Id: string;

  constructor(x: number, y: number, prevX: number, prevY: number, objectType: ObjectType, id: string) {
    this.X = x;
    this.Y = y;
    this.PrevX = prevX;
    this.PrevY = prevY;
    this.ObjectType = objectType;
    this.Id = id
  }

  public static fromJSON(data: any): Change {
    return new Change(Number(data.X), Number(data.Y), Number(data.PrevX), Number(data.PrevY), data.ObjectType as ObjectType, String(data.Id));
  }
}
