package SubmitMoveUseCase

import (
	"ChoHanJi/domain/Action"
	"ChoHanJi/domain/Player"
	"ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ErrWrongInput = errors.New("input format wrong")

type Interface interface {
	Submit(roomId Room.Id, actionType Action.Enum, msg []byte) error
}

type IActionList interface {
	SubmitMoveAction(roomId Room.Id, x, y, prevX, prevY int, id Player.Id) error
	SubmitAttackAction(roomId Room.Id, attackerId, defenderId Player.Id) error
	SubmitBonusAttackAction(roomId Room.Id, x, y int, attackerId Player.Id) error
}

type IHub interface {
	Publish(roomId, subscriberId, messageType, messageBody string) error
}

var _ IHub = (*SSEHub.SSEHub)(nil)

var _ IActionList = (*Action.List)(nil)

type Struct struct {
	al        IActionList
	hub       IHub
	validator *validator.Validate
}

func New(validator *validator.Validate, hub IHub, actionList IActionList) *Struct {
	return &Struct{actionList, hub, validator}
}

var _ Interface = (*Struct)(nil)

// Submit implements Interface.
func (s *Struct) Submit(roomId Room.Id, actionType Action.Enum, msg []byte) error {
	var id Player.Id

	switch actionType {
	case Action.Attack:
		var attackAction Action.AttackAction
		if err := json.Unmarshal(msg, &attackAction); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		// Validate
		if err := s.validator.Struct(attackAction); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		id = attackAction.AttackerId
		if err := s.al.SubmitAttackAction(roomId, attackAction.AttackerId, attackAction.DefenderId); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w", err)
		}
	case Action.Move:
		var move Action.MoveAction
		if err := json.Unmarshal(msg, &move); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		// Validate
		if err := s.validator.Struct(move); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		id = move.Id
		if err := s.al.SubmitMoveAction(roomId, move.X, move.Y, move.PrevX, move.PrevY, move.Id); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w", err)
		}
	case Action.BonusAttack:
		var action Action.BonusAttackAction
		if err := json.Unmarshal(msg, &action); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		// Validate
		if err := s.validator.Struct(action); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		id = action.Id
		if err := s.al.SubmitBonusAttackAction(roomId, action.X, action.Y, action.Id); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w", err)
		}
	case Action.Skip:
		var skip Action.SkipAction
		if err := json.Unmarshal(msg, &skip); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		// Validate
		if err := s.validator.Struct(skip); err != nil {
			return fmt.Errorf("SubmitMoveUseCase.Submit: %w: %v", ErrWrongInput, err)
		}
		id = skip.Id
	default:
		return fmt.Errorf("SubmitMoveUseCase.Submit: %w", ErrWrongInput)
	}

	if err := s.hub.Publish(string(roomId), "admin", "PlayerIsReady", string(id)); err != nil {
		return fmt.Errorf("SubmitMoveUseCase.Submit: %w", err)
	}

	return nil
}
