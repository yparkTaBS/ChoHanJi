package SubmitFightResultUseCase

import (
	"ChoHanJi/domain/Fight"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/PlayerBlocker"
	"ChoHanJi/domain/Room"
	"encoding/json"
	"fmt"

	"ChoHanJi/driven/sse/SSEHub"
	"github.com/go-playground/validator/v10"
)

type Interface interface {
	Submit(roomId Room.Id, fightId Fight.Id, submitterId Player.Id, winnerId Player.Id) error
}

type IFights interface {
	RegisterResult(roomId Room.Id, fightId Fight.Id, submitterId Player.Id, winnerId Player.Id) (*Fight.Struct, bool, error)
}

type IPlayerBlocker interface {
	Unblock(roomId Room.Id, playerId Player.Id) error
}

var _ IPlayerBlocker = (*PlayerBlocker.Struct)(nil)

type IHub interface {
	Publish(roomId, subscriberId, messageType, messageBody string) error
}

var _ IHub = (*SSEHub.Struct)(nil)
var _ IFights = (*Fight.CurrentFights)(nil)

type Struct struct {
	fights    IFights
	blocker   IPlayerBlocker
	hub       IHub
	validator *validator.Validate
}

func New(fights IFights, blocker IPlayerBlocker, hub IHub, validator *validator.Validate) *Struct {
	return &Struct{
		fights:    fights,
		blocker:   blocker,
		hub:       hub,
		validator: validator,
	}
}

var _ Interface = (*Struct)(nil)

type Request struct {
	FightId     Fight.Id  `json:"FightId" validate:"required,len=5,alphanum"`
	SubmitterId Player.Id `json:"SubmitterId" validate:"required,len=5,alphanum"`
	WinnerId    Player.Id `json:"WinnerId" validate:"required,len=5,alphanum"`
}

func (s *Struct) Submit(roomId Room.Id, fightId Fight.Id, submitterId Player.Id, winnerId Player.Id) error {
	fight, ready, err := s.fights.RegisterResult(roomId, fightId, submitterId, winnerId)
	if err != nil {
		return fmt.Errorf("SubmitFightResultUseCase.Submit: %w", err)
	}

	if !ready {
		return nil
	}

	// Unblock both participants so the processor can continue.
	_ = s.blocker.Unblock(roomId, fight.AttackerId)
	_ = s.blocker.Unblock(roomId, fight.DefenderId)

	msg, err := json.Marshal(fight)
	if err != nil {
		return fmt.Errorf("SubmitFightResultUseCase.Submit: %w", err)
	}

	if err := s.hub.Publish(string(roomId), string(fight.AttackerId), "FightResult", string(msg)); err != nil {
		return fmt.Errorf("SubmitFightResultUseCase.Submit: %w", err)
	}

	if err := s.hub.Publish(string(roomId), string(fight.DefenderId), "FightResult", string(msg)); err != nil {
		return fmt.Errorf("SubmitFightResultUseCase.Submit: %w", err)
	}

	return nil
}
