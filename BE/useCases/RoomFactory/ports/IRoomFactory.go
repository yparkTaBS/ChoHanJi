package ports

import (
	r "ChoHanJi/domain/Room"
)

type IRoomFactory interface {
	Create(width int, height int, items string) (r.Id, error)
}
