package Action

import (
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"fmt"
	"sync"
)

type List struct {
	lock sync.Mutex
	al   map[Room.Id]*actionList
}

type actionList struct {
	AttackList      []AttackAction
	MoveList        []MoveAction
	BonusAttackList []BonusAttackAction
}

func New() *List {
	return &List{
		al: make(map[Room.Id]*actionList),
	}
}

func (s *List) StartGame(roomId Room.Id) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.al[roomId]; !found {
		room := &actionList{
			make([]AttackAction, 0),
			make([]MoveAction, 0),
			make([]BonusAttackAction, 0),
		}
		s.al[roomId] = room
	}
}

func (s *List) Reset(roomId Room.Id) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return fmt.Errorf("ActionList.Reset: room not found")
	}
	room.AttackList = room.AttackList[:0]
	room.MoveList = room.MoveList[:0]
	room.BonusAttackList = room.BonusAttackList[:0]

	return nil
}

type MoveAction struct {
	X     int       `json:"X" validate:"gte=0"`
	Y     int       `json:"Y" validate:"gte=0"`
	PrevX int       `json:"PrevX" validate:"gte=0"`
	PrevY int       `json:"PrevY" validate:"gte=0"`
	Id    Player.Id `json:"Id" validate:"required,alphanum,len=5"`
}

func (s *List) SubmitMoveAction(roomId Room.Id, x, y, prevX, prevY int, id Player.Id) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return fmt.Errorf("ActionList.SubmitMoveAction: room not found")
	}

	room.MoveList = append(room.MoveList, MoveAction{x, y, prevX, prevY, id})

	return nil
}

type AttackAction struct {
	AttackerId Player.Id `json:"AttackerId" validate:"required,alphanum,len=5"`
	DefenderId Player.Id `json:"DefenderId" validate:"required,alphanum,len=5"`
}

func (s *List) SubmitAttackAction(roomId Room.Id, attackerId, defenderId Player.Id) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return fmt.Errorf("ActionList.SubmitAttackAction: room not found")
	}

	room.AttackList = append(room.AttackList, AttackAction{attackerId, defenderId})

	return nil
}

type BonusAttackAction struct {
	X  int       `json:"X" validate:"gte=0"`
	Y  int       `json:"Y" validate:"gte=0"`
	Id Player.Id `json:"Id" validate:"required,alphanum,len=5"`
}

func (s *List) SubmitBonusAttackAction(roomId Room.Id, x, y int, attackerId Player.Id) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return fmt.Errorf("ActionList.SubmitAttackAction: room not found")
	}

	room.BonusAttackList = append(room.BonusAttackList, BonusAttackAction{x, y, attackerId})

	return nil
}

type SkipAction struct {
	Id Player.Id `json:"Id" validate:"required,alphanum,len=5"`
}
