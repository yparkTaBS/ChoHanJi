package IdGenerator

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

func NewId() (string, error) {
	mapUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	s := mapUUID.String()
	s = strings.ReplaceAll(s, "-", "")

	return s[len(s)-12:], nil
}
