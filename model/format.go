package model

import (
	"errors"
	"strings"
)

var ErrParsingFormatType = errors.New("parsing input to format type")

const (
	Json Format = "json"
	Xml  Format = "xml"
	Csv  Format = "csv"
	Html Format = "html"
)

type Format string

var FormatValues []string = []string{Json.String(), Xml.String(), Csv.String(), Html.String()}

func (f Format) String() string {
	switch f {
	case Json:
		return "json"
	case Xml:
		return "xml"
	case Csv:
		return "csv"
	case Html:
		return "html"
	default:
		return "json"
	}
}

func (f *Format) Type() string {
	return "string"
}

func (f *Format) Set(input string) error {
	if !f.unmarshalText(input) && !f.unmarshalText(strings.ToLower(input)) {
		return ErrParsingFormatType
	}
	return nil
}

func (f Format) ToPath() string {
	return "/format/" + f.String()
}

func (f *Format) unmarshalText(input string) bool {
	switch input {
	case "json", "JSON":
		*f = Json
	case "xml", "XML":
		*f = Xml
	case "csv", "CSV":
		*f = Csv
	case "html", "HTML":
		*f = Html
	default:
		return false
	}
	return true
}

func ParseFormat(input string) (Format, error) {
	var f Format
	if !f.unmarshalText(input) {
		return f, ErrParsingFormatType
	}
	return f, nil
}
