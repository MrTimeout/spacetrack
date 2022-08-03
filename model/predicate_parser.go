package model

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrParsingFormatPredicate = errors.New("parsing input to predicate")

	PREDICATE_REGEXP = regexp.MustCompile(`^([a-zA-Z_]+)(.*)$`)
)

var operandsByPredicate = map[string]OperandValidator{
	PREDICATE_CCSDS_OMM_VERS:      OperandGeneralValidator{Validators: []regexp.Regexp{*CCSDS_OMM_VERS_REGEXP}, HelperText: CCSDS_OMM_VERS_HELP},
	PREDICATE_COMMENT:             OperandGeneralValidator{Validators: STRING_REGEXP_SLICE, HelperText: "Simple match of the comment field inside of the predicate"},
	PREDICATE_CREATION_DATE:       OperandGeneralValidator{Validators: DATE_REGEXP_SLICE, HelperText: "Creation date field which identifies the origin of the orbit object"},
	PREDICATE_ORIGINATOR:          OperandGeneralValidator{Validators: STRING_REGEXP_SLICE, HelperText: "Creating agency or operator (value should bespecified in an ICD)"},
	PREDICATE_OBJECT_NAME:         OperandGeneralValidator{Validators: STRING_REGEXP_SLICE, HelperText: "Spacecraft name for which the orbit state is provided. There is not specific format."},
	PREDICATE_OBJECT_ID:           OperandObjectIdValidator{HelperText: "ObjectID representing the orbital object in the format YYYY-NNNP[P] like 2000-052A"},
	PREDICATE_CENTER_NAME:         OperandGeneralValidator{Validators: STRING_REGEXP_SLICE, HelperText: "Origin of the reference name, which may be a natural solar system body. For example: EARTH, MOON, SUN..."},
	PREDICATE_REF_FRAME:           OperandExactMatchValidator{PossibleValues: REF_FRAME_VALUES, HelperText: REF_FRAME_HELP},
	PREDICATE_TIME_SYSTEM:         OperandExactMatchValidator{PossibleValues: TIME_SYSTEM_VALUES, HelperText: TIME_SYSTEM_HELP},
	PREDICATE_MEAN_ELEMENT_THEORY: OperandExactMatchValidator{PossibleValues: MEAN_ELEMENT_THEORY_VALUES, HelperText: MEAN_ELEMENT_THEORY_HELP},
	PREDICATE_EPOCH:               OperandGeneralValidator{Validators: DATE_REGEXP_SLICE, HelperText: "Epoch of state vector and optional Keplerian elements"},
	PREDICATE_MEAN_MOTION:         OperandGeneralValidator{Validators: NUMBER_REGEXP_SLICE, HelperText: "rev per day"},
	PREDICATE_ECCENTRICITY:        OperandGeneralValidator{Validators: NUMBER_REGEXP_SLICE, HelperText: "Eccentricity: https://en.wikipedia.org/wiki/Orbital_eccentricity"},
	PREDICATE_INCLINATION:         OperandGeneralValidator{Validators: NUMBER_REGEXP_SLICE, HelperText: "Inclination of the object in the orbit"},
	PREDICATE_RA_OF_ASC_NODE:      OperandGeneralValidator{Validators: NUMBER_REGEXP_SLICE, HelperText: "Right ascension of ascending node"},
	PREDICATE_ARG_OF_PERICENTER:   OperandGeneralValidator{Validators: NUMBER_REGEXP_SLICE, HelperText: "Argument of pericenter"},
	PREDICATE_COUNTRY_CODE:        OperandGeneralValidator{Validators: STRING_REGEXP_SLICE, HelperText: "Country Code"},
	PREDICATE_DECAY_DATE:          OperandGeneralValidator{Validators: DATE_REGEXP_SLICE, HelperText: "Orbital decay"},
}

func IsOperandValid(input string) bool {
	submatches := PREDICATE_REGEXP.FindStringSubmatch(input)
	if len(submatches) != 3 || submatches[2] == "" {
		return false
	}
	if operandValidator, b := operandsByPredicate[strings.ToUpper(submatches[1])]; b {
		operandSplitted := strings.Split(submatches[2], ",")
		for _, o := range operandSplitted {
			o = strings.ReplaceAll(o, "=", "")
			if !operandValidator.Validate(o) {
				return false
			}
		}
		return true
	}
	return false
}

func ToPredicates(input []string) ([]Predicate, error) {
	var err error
	result := make([]Predicate, len(input))

	for i := range input {
		if result[i], err = ToPredicate(input[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func ToPredicate(input string) (Predicate, error) {
	submatches := PREDICATE_REGEXP.FindStringSubmatch(input)
	if len(submatches) != 3 || submatches[2] == "" {
		return Predicate{}, ErrParsingFormatPredicate
	}
	return Predicate{
		Name:  submatches[1],
		Value: submatches[2],
	}, nil
}

func OperandHelp(input string) string {
	if operand, b := operandsByPredicate[strings.ToUpper(input)]; b {
		return operand.Help()
	}
	return ""
}
