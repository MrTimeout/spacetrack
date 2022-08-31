package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {

	t.Run("headers are not pair", func(t *testing.T) {
		got := Post("", []byte{}, func(r *http.Response) error { return nil }, "Content-Type")

		assert.ErrorIs(t, got, ErrHeadersMustBePair)
	})

	t.Run("new request fails when url contains invalid characters", func(t *testing.T) {
		got := Post(string([]byte{0x01, 0x02, 0x03, 0x04, 0x05}), []byte{}, func(r *http.Response) error { return nil }, "Content-Type", "application/json")

		assert.ErrorContains(t, got, "invalid control character in URL")
	})
}
