package Item

import "ChoHanJi/domain/IdGenerator"

type ItemId string

type Item struct {
	Id   ItemId
	Name string
}

var emptyItem Item

func New(name string) (Item, error) {
	strId, err := IdGenerator.NewId()
	if err != nil {
		return emptyItem, err
	}

	return Item{ItemId(strId), name}, nil
}
