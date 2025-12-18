package ports

import "ChoHanJi/domain/Map"

type IRoomFactory interface {
	Create(width int, height int, items string) (Map.Id, error)
}
