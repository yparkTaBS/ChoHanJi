package Death

import (
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"sync"
)

type List struct {
	lock sync.RWMutex
	list map[Room.Id]map[Player.Id]struct{}
}

func NewDeathList() *List {
	return &List{
		list: make(map[Room.Id]map[Player.Id]struct{}),
	}
}

func (l *List) PronounceDead(roomId Room.Id, playerId Player.Id) {
	l.lock.Lock()
	defer l.lock.Unlock()

	room, found := l.list[roomId]
	if !found {
		l.list[roomId] = make(map[Player.Id]struct{})
		room = l.list[roomId]
	}

	room[playerId] = struct{}{}
}

func (l *List) CheckIfDead(roomId Room.Id, playerId Player.Id) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	room, found := l.list[roomId]
	if !found {
		return false
	}

	_, found = room[playerId]
	return found
}

func (l *List) GetListOfDead(roomId Room.Id) []Player.Id {
	l.lock.RLock()
	defer l.lock.RUnlock()

	var deathList []Player.Id
	room, found := l.list[roomId]
	if !found {
		return deathList
	}

	for key := range room {
		deathList = append(deathList, key)
	}
	return deathList
}

func (l *List) Reset(roomId Room.Id) {
	l.lock.Lock()
	defer l.lock.Unlock()

	delete(l.list, roomId) // delete is safe even if key doesn't exist
}
