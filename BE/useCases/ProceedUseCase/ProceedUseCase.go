package ProceedUseCase

import (
	"ChoHanJi/domain/Action"
	"ChoHanJi/domain/PlayerBlocker"
	"ChoHanJi/domain/Room"
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
)

type Interface interface {
	Proceed(ctx context.Context, roomId string, logger *slog.Logger) error
}

type ActionList interface {
	GetAttackActionList(roomId Room.Id) ([]Action.AttackStruct, error)
	GetMoveActionList(roomId Room.Id) ([]Action.MoveStruct, error)
	GetBonusAttackList(roomId Room.Id) ([]Action.BonusAttackStruct, error)
	Reset(roomId Room.Id) error
}

var _ ActionList = (*Action.List)(nil)

type ActionProcessor interface {
	Process(roomId Room.Id, attacks []Action.AttackStruct, moves []Action.MoveStruct, bonusAttacks []Action.BonusAttackStruct) error
}

var _ ActionProcessor = (*Action.Processor)(nil)

type IPlayerBlocker interface {
	Initialize(roomId Room.Id)
	UnblockAllChannels(roomId Room.Id) error
}

var _ IPlayerBlocker = (*PlayerBlocker.Struct)(nil)

type Struct struct {
	al ActionList
	ap ActionProcessor
	pb IPlayerBlocker
}

var _ Interface = (*Struct)(nil)

func New(al ActionList, ap ActionProcessor, pb IPlayerBlocker) *Struct {
	return &Struct{al, ap, pb}
}

// Proceed implements Interface.
func (s *Struct) Proceed(ctx context.Context, roomId string, logger *slog.Logger) error {
	defer func() {
		_ = s.al.Reset(Room.Id(roomId))
	}()
	logger = logger.With("component", "ProceedUseCase")

	id := Room.Id(roomId)

	s.pb.Initialize(id)
	if err := s.pb.UnblockAllChannels(id); err != nil {
		return err
	}

	var totalErrors error
	attackActions, err := s.al.GetAttackActionList(id)
	totalErrors = errors.Join(totalErrors, err)
	moveActions, err := s.al.GetMoveActionList(id)
	totalErrors = errors.Join(totalErrors, err)
	bonusAttackActions, err := s.al.GetBonusAttackList(id)
	totalErrors = errors.Join(totalErrors, err)

	if totalErrors != nil {
		return totalErrors
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx,
					"Panic in ap.Process",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
			}
		}()

		if err := s.ap.Process(Room.Id(roomId), attackActions, moveActions, bonusAttackActions); err != nil {
			logger.ErrorContext(ctx, "Error Processing the Proceed Request", slog.Any("Error", err))
		}
	}()

	return nil
}
