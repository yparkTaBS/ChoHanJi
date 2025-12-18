package Player

import c "ChoHanJi/domain/Class"

type PlayerId string

type Player struct {
	Id    PlayerId
	Name  string
	Class c.Class
}
