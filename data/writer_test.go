package data

import (
	"errors"
	"testing"

	"github.com/MrTimeout/spacetrack/model"
	"github.com/stretchr/testify/assert"
)

var errDumbMarshaller = errors.New("dumb marshaller not working")

type dumbMarshaller struct{}

func (d dumbMarshaller) marshal(input any) ([]byte, error) {
	return nil, errDumbMarshaller
}

func (d dumbMarshaller) ext() string {
	return "dumb"
}

func TestGetMarshaller(t *testing.T) {
	for _, each := range []struct {
		description string
		mFormat     model.Format
		want        Marshaller
		wantErr     error
	}{
		{
			description: "xml marshaller",
			mFormat:     model.Xml,
			want:        XMLMarshaller{},
		},
		{
			description: "json marshaller",
			mFormat:     model.Json,
			want:        JSONMarshaller{},
		},
		{
			description: "csv marshaller",
			mFormat:     model.Csv,
			want:        CSVMarshaller{},
		},
		{
			description: "html marshaller",
			mFormat:     model.Html,
			want:        HTMLMarshaller{},
		},
		{
			description: "invalid marshaller",
			mFormat:     "invalid",
			wantErr:     model.ErrParsingFormatType,
		},
	} {
		t.Run(each.description, func(t *testing.T) {
			got, gotErr := getMarshaller(each.mFormat)

			assert.Equal(t, each.want, got)
			assert.ErrorIs(t, each.wantErr, gotErr)
		})
	}
}

func TestMarshal(t *testing.T) {
	type dumb struct {
		Value1 string `json:"value1" xml:"Value1" csv:"value1" html:"l=Value1,e=span"`
		Value2 int    `json:"value2" xml:"Value2" csv:"value2" html:"l=Value2,e=span"`
	}

	var d = dumb{Value1: "here some value", Value2: 2}

	for _, each := range []struct {
		description string
		m           Marshaller
		input       any
		want        string
		wantErr     error
	}{
		{
			description: "xml marshalling",
			m:           XMLMarshaller{},
			input:       d,
			want:        "<dumb><Value1>here some value</Value1><Value2>2</Value2></dumb>",
		},
		{
			description: "json marshalling",
			m:           JSONMarshaller{},
			input:       d,
			want:        "{\"value1\":\"here some value\",\"value2\":2}",
		},
		{
			description: "csv marshalling",
			m:           CSVMarshaller{},
			input:       []dumb{d},
			want:        "value1,value2\nhere some value,2\n",
		},
		{
			description: "html marshalling",
			m:           HTMLMarshaller{},
			input:       d,
			want:        "<div><span>Value1</span><span>here some value</span></div><div><span>Value2</span><span>2</span></div>",
		},
	} {
		t.Run(each.description, func(t *testing.T) {
			b, err := each.m.marshal(each.input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, each.want, string(b))
		})
	}
}

func TestMarshallerExt(t *testing.T) {
	for _, each := range []struct {
		description string
		input       Marshaller
		want        string
	}{
		{
			description: "xml ext",
			input:       XMLMarshaller{},
			want:        model.Xml.String(),
		},
		{
			description: "json ext",
			input:       JSONMarshaller{},
			want:        model.Json.String(),
		},
		{
			description: "csv ext",
			input:       CSVMarshaller{},
			want:        model.Csv.String(),
		},
		{
			description: "html ext",
			input:       HTMLMarshaller{},
			want:        model.Html.String(),
		},
	} {
		t.Run(each.description, func(t *testing.T) {
			assert.Equal(t, each.want, each.input.ext())
		})
	}
}

func TestWriter(t *testing.T) {
	t.Run("writer with marshal not working", func(t *testing.T) {
		w := WriterImpl{dumbMarshaller{}}

		gotErr := w.Write("./not_existent_file", struct{}{})

		assert.ErrorIs(t, errDumbMarshaller, gotErr)
	})
}
