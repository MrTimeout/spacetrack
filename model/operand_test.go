package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperandExactMatchValidator(t *testing.T) {
	operandTime := OperandExactMatchValidator{
		PossibleValues: TIME_SYSTEM_VALUES,
		HelperText:     TIME_SYSTEM_HELP,
	}

	t.Run("operand exact match validator success when input match even if lower case", func(t *testing.T) {
		for _, val := range operandTime.PossibleValues {
			assert.True(t, operandTime.Validate(val))
			assert.True(t, operandTime.Validate(strings.ToLower(val)))
		}
	})

	t.Run("operand extact match validator fail when input is not correct", func(t *testing.T) {
		assert.False(t, operandTime.Validate("notexistent"))
	})

	t.Run("operand exact match validator returns the correct helper text", func(t *testing.T) {
		assert.Equal(t, TIME_SYSTEM_HELP, operandTime.Help())
	})
}

func TestOperandObjectIdValidator(t *testing.T) {
	var testCases = []struct {
		description, input string
		want               bool
	}{
		{
			description: "when correct object-id is inserted",
			input:       "1960-000A",
			want:        true,
		},
		{
			description: "when correct object-id is inserted with year lower than allowed",
			input:       "1956-000A",
			want:        false,
		},
		{
			description: "when correct object-id is inserted with year higher than allowed",
			input:       "4000-000A",
			want:        false,
		},
		{
			description: "when correct object-id is inserted",
			input:       "1960-000A",
			want:        true,
		},
		{
			description: "when correct format of year is inserted",
			input:       ">1960",
			want:        true,
		},
		{
			description: "when correct format of year is inserted with year lower than allowed",
			input:       ">1950",
			want:        false,
		},
		{
			description: "when correct format of year is inserted with year higher than allowed",
			input:       "<4000",
			want:        false,
		},
		{
			description: "when not a valid input is inserted",
			input:       "notvalidinput",
			want:        false,
		},
	}
	var operandObjectIdValidator = OperandObjectIdValidator{
		HelperText: "some helper text here",
	}

	for _, tCase := range testCases {
		t.Run(tCase.description, func(t *testing.T) {
			got := operandObjectIdValidator.Validate(tCase.input)

			assert.Equal(t, tCase.want, got)
		})
	}

	t.Run("when helper text returns correctly", func(t *testing.T) {
		got := operandObjectIdValidator.Help()
		want := "some helper text here"

		assert.Equal(t, want, got)
	})

	t.Run("when input year is not a number", func(t *testing.T) {
		got := operandObjectIdValidator.isYearValid("notarealnumber")

		assert.False(t, got)
	})

}
