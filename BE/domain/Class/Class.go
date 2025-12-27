package Class

type Struct struct {
	Power         int
	Defence       int
	Range         int
	MovementSpeed int
	InitialHP     int
}

var (
	Fighter Struct = Struct{2, 1, 1, 1, 2} // can move and attack
	Ranger  Struct = Struct{2, 0, 2, 1, 2} // can attack from afar
	Rogue   Struct = Struct{2, 0, 1, 3, 2} // can move around fighter
)
