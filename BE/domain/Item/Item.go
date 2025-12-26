package Item

import (
	"ChoHanJi/domain/IdGenerator"
	"strings"
)

type Id string

type Struct struct {
	X     int `json:"X"`
	Y     int `json:"Y"`
	Id    Id
	IdStr string `json:"Id"`
	Name  string `json:"Name"`
}

func New(name string) (*Struct, error) {
	strId, err := IdGenerator.NewId()
	if err != nil {
		return nil, err
	}

	return &Struct{Id: Id(strId), IdStr: strId, Name: name}, nil
}

func DecodeItems(itemList string) []*Struct {
	var items []*Struct

	// Add items if provided
	itemList = strings.TrimSpace(itemList)
	if itemList != "" {
		itemNames := strings.SplitSeq(itemList, ",")
		for itemName := range itemNames {
			itemName = strings.TrimSpace(itemName)
			if itemName == "" {
				continue
			}

			item, _ := New(itemName)
			items = append(items, item)
		}
	}

	return items
}
