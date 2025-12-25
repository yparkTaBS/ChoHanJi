export default class Item {
  public X: number;
  public Y: number;
  public Id: string;
  public Name: string;

  constructor(x: number, y: number, id: string, name: string) {
    this.X = x;
    this.Y = y;
    this.Id = id;
    this.Name = name
  }

  public static fromJSON(data: any): Item {
    return new Item(
      Number(data.X),
      Number(data.Y),
      String(data.Id),
      String(data.Name)
    )
  }
}
