package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"os"

	"github.com/DrGrimshaw/gohtml"
	"github.com/gocarina/gocsv"
	"go.uber.org/zap"
)

// Marshaller is responsible of converting the golang struct to a string encoded with the marshal method
type Marshaller interface {
	marshal(any) ([]byte, error)
	ext() string
}

// Writer is reponsible of writing the content marshalled by the marshaller to the file passed as parameter
type Writer interface {
	Write(string, any) error
}

type WriterImpl struct {
	Marshaller
}

func (w WriterImpl) Write(file string, input any) error {
	b, err := w.marshal(input)
	if err != nil {
		return err
	}

	file = file + "." + w.ext()

	Info("writing to file", zap.String("file_name", file), zap.Int("content_size", len(b)))

	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}

	return nil
}

func getMarshaller(mFormat Format) (Marshaller, error) {
	var (
		m   Marshaller
		err error
	)

	switch mFormat {
	case Json:
		m = JSONMarshaller{}
	case Xml:
		m = XMLMarshaller{}
	case Csv:
		m = CSVMarshaller{}
	case Html:
		m = HTMLMarshaller{}
	default:
		err = ErrParsingFormatType
	}

	return m, err
}

type XMLMarshaller struct{}

func (m XMLMarshaller) marshal(input any) ([]byte, error) {
	return xml.Marshal(input)
}

func (m XMLMarshaller) ext() string {
	return Xml.String()
}

type JSONMarshaller struct{}

func (m JSONMarshaller) marshal(input any) ([]byte, error) {
	return json.Marshal(input)
}

func (m JSONMarshaller) ext() string {
	return Json.String()
}

type CSVMarshaller struct{}

func (m CSVMarshaller) marshal(input any) ([]byte, error) {
	var b bytes.Buffer

	err := gocsv.Marshal(input, &b)

	return b.Bytes(), err
}

func (m CSVMarshaller) ext() string {
	return Csv.String()
}

type HTMLMarshaller struct{}

func (h HTMLMarshaller) marshal(input any) ([]byte, error) {
	str, err := gohtml.Encode(input)
	return []byte(str), err
}

func (h HTMLMarshaller) ext() string {
	return Html.String()
}
