package model

import (
	"errors"
	"strings"
)

var ErrParsingFormatSort = errors.New("unmarshalling text to Sort type")

type Sort string

const (
	Asc  Sort = "asc"
	Desc Sort = "desc"
)

var SortValues = []string{Asc.String(), Desc.String()}

func (s Sort) String() string {
	var result = ""

	switch s {
	case Asc:
		result = "asc"
	case Desc:
		result = "desc"
	}

	return result
}

func (s *Sort) Type() string {
	return "string"
}

func (s *Sort) Set(input string) error {
	if !s.unmarshalText(input) && !s.unmarshalText(strings.ToLower(input)) {
		return ErrParsingFormatSort
	}
	return nil
}

func (s *Sort) unmarshalText(input string) bool {
	switch input {
	case "asc", "ASC":
		*s = Asc
	case "desc", "DESC":
		*s = Desc
	default:
		return false
	}
	return true
}

func ParseSort(input string) (Sort, error) {
	var s Sort
	if !s.unmarshalText(input) {
		return s, ErrParsingFormatSort
	}
	return s, nil
}
