package bootstrap

import (
	"ChoHanJi/domain/Map"
	"ChoHanJi/drivers/http/handlers"
	"ChoHanJi/drivers/http/handlers/CreateRoom"
	"ChoHanJi/useCases/RoomFactory"
	"ChoHanJi/useCases/RoomFactory/ports"
	"context"
	"fmt"

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
	); err != nil {
		return err
	}

	return nil
}

func RegisterUseCases(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		RoomFactory.NewRoomFactory,
		o.As[ports.IRoomFactory],
		o.AsSingleton,
	); err != nil {
		return err
	}
	return nil
}

func RegisterDomains(ctx context.Context, builder *cb.ContainerBuilder) error {
	if err := builder.Register(
		Map.New,
		o.AsSingleton,
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
