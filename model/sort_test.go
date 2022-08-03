package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sortTestCases = []struct {
	sort Sort
	str  string
}{
	{
		sort: Asc,
		str:  "asc",
	},
	{
		sort: Desc,
		str:  "desc",
	},
}

func TestSort(t *testing.T) {
	t.Run("", func(t *testing.T) {
		for _, tCase := range sortTestCases {
			assert.Equal(t, tCase.str, tCase.sort.String())
		}
	})

	t.Run("format set returns nil when text exists", func(t *testing.T) {
		for _, tCase := range sortTestCases {
			got, err := ParseSort(tCase.str)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tCase.sort, got)
		}
	})

	t.Run("format set returns error when text don't exist", func(t *testing.T) {
		_, got := ParseSort("notexistent")
		want := ErrParsingFormatSort

		assert.ErrorIs(t, want, got)
	})

	t.Run("type of sort must return string", func(t *testing.T) {
		for _, tCase := range sortTestCases {
			assert.Equal(t, "string", tCase.sort.Type())
		}
	})

	t.Run("set correctly when input is correct", func(t *testing.T) {
		for _, tCase := range sortTestCases {
			var sort Sort

			got := sort.Set(tCase.str)

			assert.Nil(t, got)
		}
	})

	t.Run("set returns error when input is incorrect", func(t *testing.T) {
		var sort Sort

		got := sort.Set("notexistent")
		want := ErrParsingFormatSort

		assert.ErrorIs(t, want, got)
	})
}
