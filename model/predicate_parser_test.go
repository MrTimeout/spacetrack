package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOperandValid(t *testing.T) {
	var testCases = []struct {
		description, input string
		want               bool
	}{
		{
			description: "submatches are less than the required 3",
			input:       "incorrectpredicate",
			want:        false,
		},
		{
			description: "predicate is not present in the allowed ones",
			input:       "notexistent=single comment",
			want:        false,
		},
		{
			description: "predicate is present, but the filter doesn't match the criteria",
			input:       "comment>singlet comment",
			want:        false,
		},
		{
			description: "predicate is present and has only one filter",
			input:       "comment=single comment",
			want:        true,
		},
		{
			description: "predicate is present and has more than one filter",
			input:       "comment=single comment,^another one,~~hey there",
			want:        true,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.description, func(t *testing.T) {
			got := IsOperandValid(tCase.input)

			assert.Equal(t, tCase.want, got)
		})
	}
}

func TestToPredicate(t *testing.T) {
	var testCases = []struct {
		description string
		input       []string
		want        []Predicate
		errWanted   error
	}{
		{
			description: "predicate with an incorrect format, because is has less than 3 matches",
			input:       []string{"notarealpredicate"},
			want:        nil,
			errWanted:   ErrParsingFormatPredicate,
		},
		{
			description: "predicates with real predicates inside",
			input:       []string{"comment=single comment", "object_id=1960-000A"},
			want: []Predicate{
				{
					Name:  "comment",
					Value: "=single comment",
				},
				{
					Name:  "object_id",
					Value: "=1960-000A",
				},
			},
			errWanted: nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.description, func(t *testing.T) {
			got, gotErr := ToPredicates(tCase.input)

			assert.Equal(t, tCase.want, got)
			assert.ErrorIs(t, tCase.errWanted, gotErr)
		})
	}
}

func TestOperandHelp(t *testing.T) {
	var testCases = []struct {
		description, input, want string
	}{
		{
			description: "input operand not existent will return empty string",
			input:       "notexistent",
			want:        "",
		},
		{
			description: "input operand existent in upper case will return successfully",
			input:       "COMMENT",
			want:        "Simple match of the comment field inside of the predicate",
		},
		{
			description: "input operand existent in lower case will return successfully",
			input:       "comment",
			want:        "Simple match of the comment field inside of the predicate",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.description, func(t *testing.T) {
			got := OperandHelp(tCase.input)

			assert.Equal(t, tCase.want, got)
		})
	}
}
