package Fight

import (
	"ChoHanJi/domain/Game"
	"ChoHanJi/domain/IdGenerator"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
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
			return &Struct{Id(id), gameType, attId, nil, defId, nil}, nil
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
