package Fight

import (
	"ChoHanJi/domain/Game"
	"ChoHanJi/domain/IdGenerator"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"fmt"
	"sync"
)

type Id string

type Struct struct {
	Id             Id
	Type           Game.Type
	AttackerId     Player.Id
	AttackerResult any
	DefenderId     Player.Id
	DefenderResult any
	WinnerId       Player.Id `json:"WinnerId,omitempty"`
	submissions    map[Player.Id]struct{}
	resolved       bool
}

type CurrentFights struct {
	lock          sync.RWMutex
	currentFights map[Room.Id]map[Id]*Struct
}

func (cf *CurrentFights) newFight(room map[Id]*Struct, gameType Game.Type, attId, defId Player.Id) (*Struct, error) {
	for {
		id, err := IdGenerator.NewId()
		if err != nil {
			return nil, err
		}

		_, found := room[Id(id)]
		if !found {
			return &Struct{
				Id:          Id(id),
				Type:        gameType,
				AttackerId:  attId,
				DefenderId:  defId,
				submissions: make(map[Player.Id]struct{}),
			}, nil
		}
	}
}

func New() *CurrentFights {
	return &CurrentFights{
		currentFights: make(map[Room.Id]map[Id]*Struct),
	}
}

func (cf *CurrentFights) Create(roomId Room.Id, gameType Game.Type, attId, defId Player.Id) (*Struct, error) {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	if _, found := cf.currentFights[roomId]; !found {
		cf.currentFights[roomId] = make(map[Id]*Struct)
	}

	room := cf.currentFights[roomId]
	fight, err := cf.newFight(room, gameType, attId, defId)
	if err != nil {
		return nil, err
	}

	room[fight.Id] = fight

	return fight, nil
}

func (cf *CurrentFights) RegisterResult(roomId Room.Id, fightId Id, submitterId Player.Id, winnerId Player.Id) (*Struct, bool, error) {
	cf.lock.Lock()
	defer cf.lock.Unlock()

	room, found := cf.currentFights[roomId]
	if !found {
		return nil, false, fmt.Errorf("CurrentFights.RegisterResult: room not found")
	}

	fight, found := room[fightId]
	if !found {
		return nil, false, fmt.Errorf("CurrentFights.RegisterResult: fight not found")
	}

	if fight.resolved {
		return fight, false, nil
	}

	if submitterId != fight.AttackerId && submitterId != fight.DefenderId {
		return nil, false, fmt.Errorf("CurrentFights.RegisterResult: player not part of fight")
	}

	if winnerId != fight.AttackerId && winnerId != fight.DefenderId {
		return nil, false, fmt.Errorf("CurrentFights.RegisterResult: winner not part of fight")
	}

	fight.submissions[submitterId] = struct{}{}

	if fight.WinnerId == "" {
		fight.WinnerId = winnerId
	}

	ready := len(fight.submissions) >= 2
	if ready {
		fight.resolved = true
	}

	return fight, ready, nil
}
