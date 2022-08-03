package model

import (
	"strconv"
)

type Limit struct {
	Max  int
	Skip int
}

func (l Limit) ToPath() string {
	if l.Max <= 0 {
		return ""
	}
	var path = "/limit/" + strconv.Itoa(l.Max)
	if l.Skip > 0 {
		path += "," + strconv.Itoa(l.Skip)
	}
	return path
}
