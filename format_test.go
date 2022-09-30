package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var formatTestCases = []struct {
	format Format
	str    string
	path   string
}{
	{
		format: Json,
		str:    "json",
		path:   "/format/json",
	},
	{
		format: Html,
		str:    "html",
		path:   "/format/html",
	},
	{
		format: Csv,
		str:    "csv",
		path:   "/format/csv",
	},
	{
		format: Xml,
		str:    "xml",
		path:   "/format/xml",
	},
}

func TestFormat(t *testing.T) {
	t.Run("format set returns nil when text exists", func(t *testing.T) {
		for _, f := range formatTestCases {
			var format Format

			assert.Nil(t, format.Set(f.str))
		}
	})

	t.Run("format set returns error when text don't exist", func(t *testing.T) {
		var format Format

		got := format.Set("notexistent")
		want := ErrParsingFormatType

		assert.ErrorIs(t, got, want)
	})

	t.Run("parse format returns the right format when exist", func(t *testing.T) {
		for _, want := range formatTestCases {
			got, err := ParseFormat(want.str)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, want.format, got)
		}
	})

	t.Run("parse format returns error when not exist", func(t *testing.T) {
		_, got := ParseFormat("notexistent")
		want := ErrParsingFormatType

		assert.ErrorIs(t, got, want)
	})

	t.Run("to path returns format of the Format type", func(t *testing.T) {
		for _, want := range formatTestCases {
			got := want.format.ToPath()

			assert.Equal(t, got, want.path)
		}
	})

	t.Run("format type is string", func(t *testing.T) {
		for _, f := range formatTestCases {
			assert.Equal(t, "string", f.format.Type())
		}
	})
}
