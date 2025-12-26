package bootstrap

import (
	"ChoHanJi/domain/Room"
	"ChoHanJi/driven/sse/SSEHub"
	"ChoHanJi/drivers/http/handlers"
	"ChoHanJi/drivers/http/handlers/CreateCharacter"
	"ChoHanJi/drivers/http/handlers/CreateRoom"
	AdminGameStatus "ChoHanJi/drivers/http/handlers/GameStatus/Admin"
	"ChoHanJi/drivers/http/handlers/PlayerRoom"
	"ChoHanJi/drivers/http/handlers/StartGame"
	"ChoHanJi/drivers/http/handlers/WaitingRoom"
	"ChoHanJi/useCases/AdminWaitingRoomUseCase"
	"ChoHanJi/useCases/CharacterFactory"
	"ChoHanJi/useCases/GameStatus"
	"ChoHanJi/useCases/PlayerWaitingRoomUseCase"
	"ChoHanJi/useCases/RoomFactory"
	"ChoHanJi/useCases/RoomFactory/ports"
	"ChoHanJi/useCases/StartGameUseCase"
	"context"
	"fmt"
	"net/http"

	gi "github.com/TaBSRest/GoFac/interfaces"
	cb "github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
	"github.com/go-playground/validator/v10"

	o "github.com/TaBSRest/GoFac/pkg/Options/Registration"
)

func Register(ctx context.Context) gi.Container {
	cb := cb.New()

	if err := RegisterDrivers(ctx, cb); err != nil {
		panic(fmt.Errorf("could not register drivers! %w", err))
	}

	if err := RegisterUseCases(ctx, cb); err != nil {
		panic(fmt.Errorf("could not register use cases! %w", err))
	}

	if err := RegisterDomains(ctx, cb); err != nil {
		panic(fmt.Errorf("could not register domains! %w", err))
	}

	if err := RegisterDriven(ctx, cb); err != nil {
		panic(fmt.Errorf("could not register driven adapters! %w", err))
	}

	if err := RegisterExternalDependencies(ctx, cb); err != nil {
		panic(fmt.Errorf("could not register external dependencies! %w", err))
	}

	container, err := cb.Build()
	if err != nil {
		panic(fmt.Errorf("could not build the container builder! %w", err))
	}

	return container
}

func RegisterDrivers(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		CreateRoom.New,
		o.AsSingleton,
		o.Named(string(handlers.POSTRoom)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	if err := builder.Register(
		WaitingRoom.New,
		o.AsSingleton,
		o.Named(string(handlers.GETRoomAdmin)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	if err := builder.Register(
		CreateCharacter.New,
		o.AsSingleton,
		o.Named(string(handlers.POSTCharacter)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	if err := builder.Register(
		PlayerRoom.New,
		o.AsSingleton,
		o.Named(string(handlers.GETPlayerEvent)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	if err := builder.Register(
		StartGame.New,
		o.AsSingleton,
		o.Named(string(handlers.POSTGameStart)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	if err := builder.Register(
		AdminGameStatus.New,
		o.AsSingleton,
		o.Named(string(handlers.GETAdminGameStatus)),
		o.As[http.Handler],
	); err != nil {
		return err
	}

	return nil
}

func RegisterUseCases(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		RoomFactory.NewRoomFactory,
		o.As[ports.UseCaseInterface],
		o.AsSingleton,
	); err != nil {
		return err
	}

	if err := builder.Register(
		AdminWaitingRoomUseCase.New,
		o.AsSingleton,
		o.As[AdminWaitingRoomUseCase.UseCaseInterface],
	); err != nil {
		return err
	}

	if err := builder.Register(
		CharacterFactory.New,
		o.AsSingleton,
		o.As[CharacterFactory.UseCaseInterface],
	); err != nil {
		return err
	}

	if err := builder.Register(
		PlayerWaitingRoomUseCase.New,
		o.AsSingleton,
		o.As[PlayerWaitingRoomUseCase.UseCaseInterface],
	); err != nil {
		return err
	}

	if err := builder.Register(
		StartGameUseCase.New,
		o.AsSingleton,
		o.As[StartGameUseCase.UseCaseInterface],
	); err != nil {
		return err
	}

	if err := builder.Register(
		GameStatus.New,
		o.AsSingleton,
		o.As[GameStatus.Interface],
	); err != nil {
		return err
	}

	return nil
}

func RegisterDomains(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		Room.New,
		o.AsSingleton,
	); err != nil {
		return err
	}
	return nil
}

func RegisterDriven(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		SSEHub.New,
		o.AsSingleton,
		o.As[AdminWaitingRoomUseCase.IHub],
		o.As[PlayerWaitingRoomUseCase.IHub],
		o.As[StartGameUseCase.IHub],
	); err != nil {
		return err
	}

	if err := builder.Register(
		SSEHub.New,
		o.AsSingleton,
		o.As[GameStatus.IHub],
	); err != nil {
		return err
	}

	return nil
}

func RegisterExternalDependencies(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		func() *validator.Validate {
			return validator.New()
		},
		o.AsSingleton,
	); err != nil {
		return err
	}
	return nil
}
