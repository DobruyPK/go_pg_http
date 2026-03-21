package user

import "strings"

type Name struct {
	value string
}

func NewName(raw string) (Name, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Name{}, ErrInvalidUserName
	}
	return Name{value: trimmed}, nil
}

func (n Name) String() string {
	return n.value
}
