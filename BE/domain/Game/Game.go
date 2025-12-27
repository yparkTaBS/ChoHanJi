package Game

type Type string

const (
	EvenOrOdd        Type = "EvenOrOdd"
	BiggerDice       Type = "BiggerDice"
	GuessNumber      Type = "GuessNumber"
	RockPaperScissor Type = "RockPaperScissor"
	Station1         Type = "Station1"
	Station2         Type = "Station2"
	Station3         Type = "Station3"
)

var List = []Type{
	EvenOrOdd,
	BiggerDice,
	GuessNumber,
	RockPaperScissor,
	Station1,
	Station2,
	Station3,
}
