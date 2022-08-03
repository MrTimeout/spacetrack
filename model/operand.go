package model

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	NUMBER_REGEXP       = regexp.MustCompile(`^(\^|~~|<|>)?\d+(\.\d+)?$`)
	NUMBER_REGEXP_RANGE = regexp.MustCompile(`^\d+(\.\d+)?--\d+(\.\d+)?$`)

	STRING_REGEXP = regexp.MustCompile(`^(\^|~~)?[\p{L}_ -/]+$`)

	DATE_FORMAT      = "2006-01-02"
	DATE_TIME_FORMAT = "2006-01-02 16:04:05"

	DATE_FORMAT_REGEXP      = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	DATE_TIME_FORMAT_REGEXP = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)

	DATE_REGEXP       = regexp.MustCompile(`^(>|<)?now(-\d+(\.\d+)?)?$`)
	DATE_REGEXP_RANGE = regexp.MustCompile(`^now(-\d+(\.\d+)?)?--now(-\d+(\.\d+)?)?$`)
	NULL_VALUE_REGEXP = regexp.MustCompile(`^(<>)?null-val$`)

	CCSDS_OMM_VERS_REGEXP = regexp.MustCompile(`^[0-2]\.\d$`)
	CCSDS_OMM_VERS_HELP   = "Version number of the document, it has to be of the format x.y"

	OBJECT_ID_REGEXP      = regexp.MustCompile(`^(\^|~~)?(\d{4})-\d{3}[A-Z]{1,2}$`)
	OBJECT_ID_YEAR_REGEXP = regexp.MustCompile(`^(\^|~~|<|>)?(\d+)`)

	NUMBER_REGEXP_SLICE = []regexp.Regexp{*NUMBER_REGEXP, *NUMBER_REGEXP_RANGE}
	STRING_REGEXP_SLICE = []regexp.Regexp{*STRING_REGEXP}
	DATE_REGEXP_SLICE   = []regexp.Regexp{
		*DATE_FORMAT_REGEXP,
		*DATE_TIME_FORMAT_REGEXP,
		*DATE_REGEXP,
		*DATE_REGEXP_RANGE,
		*NULL_VALUE_REGEXP,
	}

	MEAN_ELEMENT_THEORY_VALUES = []string{"SGP4", "DSST", "USM"}
	MEAN_ELEMENT_THEORY_HELP   = `Description of the mean element theory. Indicates the proper method to employ to propagate the state.`

	REF_FRAME_VALUES = []string{"EME2000", "GCRF", "ICRF", "ITRF2000", "ITRF-93", "ITRF-97", "MCI", "TDR", "TEME", "TOD"}
	REF_FRAME_HELP   = `EME2000: Earth Mean Equator and Equinox of J2000
	GCRF: Geocentric Celestial Reference Frame
	GRC: Greenwich Rotating Coordinates
	ICRF: International Celestial Reference Frame
	ITRF2000: International Terrestrial Reference Frame 2000
	ITRF-93: International Terrestrial Reference Frame 1993
	ITRF-97: International Terrestrial Reference Frame 1997
	MCI: Mars Centered Inertial
	TDR: True of Date, Rotating
	TEME: True Equator Mean Equinox (see below)
	TOD: True of Date`

	TIME_SYSTEM_VALUES = []string{"UTC", "TAI", "TT", "GPS", "TDB", "TCB"}
	TIME_SYSTEM_HELP   = `Possible values are:
	UTC: Universal Coordinated Time
	TAI: Internation Atomic Time
	TT: Terrestrial Time
	GMST: Greenwich Mean Sidereal Time
	GPS: GPS Control Segment
	MET: Mission Elapsed Time
	MRT: Mission Relative Time
	SCLK: Spacecraft Clock (receiver)
	TGG: Geocentric Coordinate Time
	TDB: Barycentric Dynamical Time
	TCB: Barycentric Coordinate Time
	UT1: Universal Time`
)

type OperandValidator interface {
	Validate(input string) bool
	Help() string
}

type OperandGeneralValidator struct {
	Validators []regexp.Regexp
	HelperText string
}

func (o OperandGeneralValidator) Validate(input string) bool {
	for _, regex := range o.Validators {
		if regex.MatchString(input) {
			return true
		}
	}
	return false
}

func (o OperandGeneralValidator) Help() string {
	return o.HelperText
}

type OperandExactMatchValidator struct {
	PossibleValues []string
	HelperText     string
}

func (o OperandExactMatchValidator) Validate(input string) bool {
	il := strings.ToUpper(input)
	for _, pv := range o.PossibleValues {
		if il == pv {
			return true
		}
	}
	return false
}

func (o OperandExactMatchValidator) Help() string {
	return o.HelperText
}

type OperandObjectIdValidator struct {
	HelperText string
}

func (o OperandObjectIdValidator) Validate(input string) bool {
	if submatches := OBJECT_ID_REGEXP.FindStringSubmatch(input); len(submatches) > 1 && strings.Contains(submatches[0], "-") {
		return o.isYearValid(submatches[len(submatches)-1])
	}
	if result := OBJECT_ID_YEAR_REGEXP.FindStringSubmatch(input); len(result) > 1 {
		return o.isYearValid(result[len(result)-1])
	}
	return false
}

func (o OperandObjectIdValidator) isYearValid(input string) bool {
	if year, err := strconv.Atoi(input); err == nil {
		return year <= time.Now().Year() && year >= 1957
	}
	return false
}

func (o OperandObjectIdValidator) Help() string {
	return o.HelperText
}
