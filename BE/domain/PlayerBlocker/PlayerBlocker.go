package PlayerBlocker

import (
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyBlocked = errors.New("the played is already blocked")
)

type Struct struct {
	lock          sync.RWMutex
	actionBlocker map[Room.Id]map[Player.Id]chan struct{}
}

func New() *Struct {
	blocker := make(map[Room.Id]map[Player.Id]chan struct{})
	return &Struct{sync.RWMutex{}, blocker}
}

func (s *Struct) Initialize(roomId Room.Id) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, found := s.actionBlocker[roomId]; !found {
		s.actionBlocker[roomId] = make(map[Player.Id]chan struct{})
	}
}

func (s *Struct) Block(roomId Room.Id, playerId Player.Id) error {
	location := "PlayerBlocker.Block"

	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.actionBlocker[roomId]
	if !found {
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	_, found = room[playerId]
	if found {
		return fmt.Errorf("%s: %w", location, ErrAlreadyBlocked)
	}

	room[playerId] = make(chan struct{})
	return nil
}

func (s *Struct) BlockPair(roomId Room.Id, p1, p2 Player.Id) error {
	location := "PlayerBlocker.BlockPair"

	s.lock.Lock()
	defer s.lock.Unlock()

	room, found := s.actionBlocker[roomId]
	if !found {
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	if _, ok := room[p1]; ok {
		return fmt.Errorf("%s: %w", location, ErrAlreadyBlocked)
	}
	if _, ok := room[p2]; ok {
		return fmt.Errorf("%s: %w", location, ErrAlreadyBlocked)
	}

	room[p1] = make(chan struct{})
	room[p2] = make(chan struct{})
	return nil
}

func (s *Struct) Unblock(roomId Room.Id, playerId Player.Id) error {
	location := "PlayerBlocker.Unblock"

	s.lock.Lock()
	room, found := s.actionBlocker[roomId]
	if !found {
		s.lock.Unlock()
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	player, found := room[playerId]
	if !found {
		s.lock.Unlock()
		return fmt.Errorf("%s: Player %w", location, ErrNotFound)
	}
	delete(s.actionBlocker[roomId], playerId)
	s.lock.Unlock()

	close(player)

	return nil
}

func (s *Struct) UnblockAllChannels(roomId Room.Id) error {
	location := "PlayerBlocker.UnblockAllChannels"

	s.lock.RLock()
	room, found := s.actionBlocker[roomId]
	if !found {
		s.lock.RUnlock()
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	playerIds := make([]Player.Id, 0, len(room))
	for pid := range room {
		playerIds = append(playerIds, pid)
	}
	s.lock.RUnlock()

	for _, pid := range playerIds {
		_ = s.Unblock(roomId, pid)
	}

	return nil
}

func (s *Struct) WaitUntilUnblocked(roomId Room.Id, playerId Player.Id) error {
	location := "PlayerBlocker.WaitUntilUnblocked"

	s.lock.RLock()

	room, found := s.actionBlocker[roomId]
	if !found {
		s.lock.RUnlock()
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	player, found := room[playerId]
	s.lock.RUnlock()

	if !found {
		return nil
	}

	<-player

	return nil
}

func (s *Struct) WaitUntilAllAreUnblocked(roomId Room.Id) error {
	location := "PlayerBlocker.WaitUntilAllAreUnblocked"

	s.lock.RLock()

	room, found := s.actionBlocker[roomId]
	if !found {
		s.lock.RUnlock()
		return fmt.Errorf("%s: Room %w", location, ErrNotFound)
	}

	var players []chan struct{}
	for _, ch := range room {
		players = append(players, ch)
	}
	s.lock.RUnlock()

	for _, ch := range players {
		<-ch
	}

	return nil
}
