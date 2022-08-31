package data

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"os"

	"github.com/DrGrimshaw/gohtml"
	"github.com/MrTimeout/spacetrack/model"
	l "github.com/MrTimeout/spacetrack/utils"
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

	l.Info("writing to file", zap.String("file_name", file), zap.Int("content_size", len(b)))

	if err := os.WriteFile(file, b, 0644); err != nil {
		return err
	}

	return nil
}

func getMarshaller(mFormat model.Format) (Marshaller, error) {
	var (
		m   Marshaller
		err error
	)

	switch mFormat {
	case model.Json:
		m = JSONMarshaller{}
	case model.Xml:
		m = XMLMarshaller{}
	case model.Csv:
		m = CSVMarshaller{}
	case model.Html:
		m = HTMLMarshaller{}
	default:
		err = model.ErrParsingFormatType
	}

	return m, err
}

type XMLMarshaller struct{}

func (m XMLMarshaller) marshal(input any) ([]byte, error) {
	return xml.Marshal(input)
}

func (m XMLMarshaller) ext() string {
	return model.Xml.String()
}

type JSONMarshaller struct{}

func (m JSONMarshaller) marshal(input any) ([]byte, error) {
	return json.Marshal(input)
}

func (m JSONMarshaller) ext() string {
	return model.Json.String()
}

type CSVMarshaller struct{}

func (m CSVMarshaller) marshal(input any) ([]byte, error) {
	var b bytes.Buffer

	err := gocsv.Marshal(input, &b)

	return b.Bytes(), err
}

func (m CSVMarshaller) ext() string {
	return model.Csv.String()
}

type HTMLMarshaller struct{}

func (h HTMLMarshaller) marshal(input any) ([]byte, error) {
	str, err := gohtml.Encode(input)
	return []byte(str), err
}

func (h HTMLMarshaller) ext() string {
	return model.Html.String()
}
