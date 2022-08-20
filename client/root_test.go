package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/MrTimeout/spacetrack/utils"
	"github.com/stretchr/testify/assert"
)

var ErrDumbReadCloser = errors.New("dumb read closer")

type DumbReadCloser struct{}

func (DumbReadCloser) Read(p []byte) (n int, err error) {
	return 0, ErrDumbReadCloser
}

func (DumbReadCloser) Close() error {
	return nil
}

func TestReadResponse(t *testing.T) {
	t.Run("incorrect content type", func(t *testing.T) {
		_, got := readResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("I am a reader")),
			Header: map[string][]string{
				"Content-Type": {"application/xml"},
			},
		})

		assert.ErrorIs(t, got, ErrIncorrectContentType)
	})

	t.Run("body is not readed correctly", func(t *testing.T) {
		_, got := readResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       DumbReadCloser{},
		})

		assert.ErrorIs(t, got, ErrDumbReadCloser)
	})

	t.Run("header is application/json but body is not", func(t *testing.T) {
		_, got := readResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("this is not a json string")),
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		})

		assert.NotNil(t, got)
	})

	t.Run("body is readed correctly and expected value is returned", func(t *testing.T) {
		want := []SpaceOrbitalObj{
			{Comment: "Hello there"},
			{Comment: "Another comment"},
		}
		b, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}

		got, err := readResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(b)),
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.ElementsMatch(t, want, got)
	})
}

func TestBuildURL(t *testing.T) {
	t.Run("build URL successfully", func(t *testing.T) {
		var want = "https://www.space-track.org/basicspacedata/query/class/gp"
		got, err := buildURL()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, want, got)
	})
}

func TestPostAuth(t *testing.T) {
	t.Run("body is not readed correctly", func(t *testing.T) {
		got := postAuth(nil)(&http.Response{
			StatusCode: http.StatusOK,
			Body:       DumbReadCloser{},
		})

		assert.ErrorIs(t, got, ErrDumbReadCloser)
	})
}

func TestFetchData(t *testing.T) {
	utils.Configure(utils.Logger{})
	sta := utils.SpaceTrackAuth{
		Identity: "identity",
		Password: "password",
		Cookie: &http.Cookie{
			Name:    "chocolatechip",
			Value:   "chocolatechip=u4cqq2pma6k6kjhedvt6pp3es1bfmndh",
			Expires: time.Now().Add(time.Hour),
		},
	}

	t.Run("invalid characters when passing one that is not allowed", func(t *testing.T) {
		_, got := FetchData(&sta, string([]byte{0x01, 0x02, 0x03, 0x04, 0x05}), true)

		assert.ErrorContains(t, got, "invalid control character in URL")
	})
}
