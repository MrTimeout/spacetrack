package main

import (
	"errors"
	"strings"
)

var ErrParsingRestCall = errors.New("parsing input to rest call")

const (
	Tle   RestCall = "tle"
	Cdm   RestCall = "cdm"
	Decay RestCall = "dec"
	All   RestCall = "all"
)

type RestCall string

var RestCallValues []string = []string{Tle.String(), Cdm.String(), Decay.String(), All.String()}

func (rc RestCall) String() string {
	var result = "all"

	switch rc {
	case Tle:
		result = "tle"
	case Cdm:
		result = "cdm"
	case Decay:
		result = "dec"
	case All:
		result = "all"
	}

	return result
}

func (rc RestCall) Type() string {
	return "string"
}

func (rc *RestCall) Set(input string) error {
	if !rc.unmarshalText(input) && !rc.unmarshalText(strings.ToLower(input)) {
		return ErrParsingRestCall
	}
	return nil
}

func (rc *RestCall) unmarshalText(input string) bool {
	switch input {
	case "tle", "TLE":
		*rc = Tle
	case "cdm", "CDM":
		*rc = Cdm
	case "dec", "DEC":
		*rc = Decay
	case "all", "ALL":
		*rc = All
	default:
		return false
	}

	return true
}

func parseRestCall(input string) (RestCall, error) {
	var rc RestCall
	if !rc.unmarshalText(input) {
		return rc, ErrParsingRestCall
	}
	return rc, nil
}
