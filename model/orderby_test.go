package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var orderbyTestCases = []struct {
	description string
	got         OrderBy
	want        string
}{
	{
		description: "when value is not present inside the possible By values, it retursn empty",
		got:         OrderBy{By: "notexistent", Sort: "asc"},
		want:        "",
	},
	{
		description: "when value is present and sort is not empty",
		got:         OrderBy{By: "FILE", Sort: "asc"},
		want:        "/orderby/FILE asc",
	},
	{
		description: "when value is present with lowercase and sort is not empty",
		got:         OrderBy{By: "file", Sort: "desc"},
		want:        "/orderby/file desc",
	},
}

func TestOrderBy(t *testing.T) {
	for _, tCase := range orderbyTestCases {
		t.Run(tCase.description, func(t *testing.T) {
			assert.Equal(t, tCase.want, tCase.got.ToPath())
		})
	}
}
