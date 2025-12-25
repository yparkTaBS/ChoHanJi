export default class Item {
  public Id: string;
  public Name: string;

  constructor(id: string, name: string) {
    this.Id = id;
    this.Name = name
  }

  public static fromJSON(data: any): Item {
    return new Item(
      String(data.Id),
      String(data.Name)
    )
  }
}
