package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var restCallTestCases = []struct {
	restCall RestCall
	str      string
}{
	{
		restCall: Tle,
		str:      "tle",
	},
	{
		restCall: Cdm,
		str:      "cdm",
	},
	{
		restCall: Decay,
		str:      "dec",
	},
	{
		restCall: All,
		str:      "all",
	},
}

func TestRestCall(t *testing.T) {
	t.Run("restCall set returns nil when text exists", func(t *testing.T) {
		for _, rc := range restCallTestCases {
			var restCall RestCall

			assert.Nil(t, restCall.Set(rc.str))
		}
	})

	t.Run("restCall set returns error when text don't exist", func(t *testing.T) {
		var restCall RestCall

		got := restCall.Set("notexistent")
		want := ErrParsingRestCall

		assert.ErrorIs(t, got, want)
	})

	t.Run("parse restCall returns the right restCall when exist", func(t *testing.T) {
		for _, want := range restCallTestCases {
			got, err := parseRestCall(want.str)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, want.restCall, got)
		}
	})

	t.Run("parse restCall returns error when not exist", func(t *testing.T) {
		_, got := parseRestCall("notexistent")
		want := ErrParsingRestCall

		assert.ErrorIs(t, got, want)
	})

	t.Run("restCall type is string", func(t *testing.T) {
		for _, f := range restCallTestCases {
			assert.Equal(t, "string", f.restCall.Type())
		}
	})
}
