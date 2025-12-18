package Class

type Class struct {
	Power         int
	Defence       int
	Range         int
	MovementSpeed int
	InitialHP     int
}

var (
	Fighter Class = Class{2, 1, 1, 1, 2} // can move and attack
	Ranger  Class = Class{2, 0, 2, 1, 2} // can attack from afar
	Rogue   Class = Class{2, 0, 1, 2, 2} // can move around fighter
)
