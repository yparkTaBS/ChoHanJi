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
	AttackList      []AttackStruct
	MoveList        []MoveStruct
	BonusAttackList []BonusAttackStruct
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
			make([]AttackStruct, 0),
			make([]MoveStruct, 0),
			make([]BonusAttackStruct, 0),
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

func (s *List) GetAttackActionList(roomId Room.Id) ([]AttackStruct, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return nil, fmt.Errorf("ActionList.GetAttackActionList: room not found")
	}

	return room.AttackList, nil
}

func (s *List) GetMoveActionList(roomId Room.Id) ([]MoveStruct, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return nil, fmt.Errorf("ActionList.GetMoveActionList: room not found")
	}

	return room.MoveList, nil
}

func (s *List) GetBonusAttackList(roomId Room.Id) ([]BonusAttackStruct, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.al[roomId]
	if !found {
		return nil, fmt.Errorf("ActionList.GetBonusAttackList: room not found")
	}

	return room.BonusAttackList, nil
}

type MoveStruct struct {
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

	room.MoveList = append(room.MoveList, MoveStruct{x, y, prevX, prevY, id})

	return nil
}

type AttackStruct struct {
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

	room.AttackList = append(room.AttackList, AttackStruct{attackerId, defenderId})

	return nil
}

type BonusAttackStruct struct {
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

	room.BonusAttackList = append(room.BonusAttackList, BonusAttackStruct{x, y, attackerId})

	return nil
}

type SkipStruct struct {
	Id Player.Id `json:"Id" validate:"required,alphanum,len=5"`
}
