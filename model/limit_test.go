package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var limitTestCases = []struct {
	description string
	limit       Limit
	want        string
}{
	{
		description: "when limit is 0 or less, empty string is returned",
		limit:       Limit{Max: 0, Skip: 2},
		want:        "",
	},
	{
		description: "when limit is more than 0 and skip is 0 or less",
		limit:       Limit{Max: 5, Skip: 0},
		want:        "/limit/5",
	},
	{
		description: "when limit is more than 0 and skip is more than 0 also",
		limit:       Limit{Max: 5, Skip: 5},
		want:        "/limit/5,5",
	},
}

func TestLimit(t *testing.T) {
	for _, testCase := range limitTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			got := testCase.limit.ToPath()

			assert.Equal(t, testCase.want, got)
		})
	}
}
