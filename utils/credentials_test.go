package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	plaintext = "Hello world"
	key       = "12345678901234567890123456789012"
)

var ciphertext = "yALnSjIiOiljHpQ="

func TestEncrypt(t *testing.T) {
	t.Run("encrypt plaintext with 32 characters key value", func(t *testing.T) {
		got, err := Encrypt([]byte(plaintext), []byte(key))
		want := ciphertext
		if err != nil {
			t.Error(err)
		}

		t.Log(string(got))

		assert.Equal(t, string(want), string(got))
	})

	t.Run("encrypt plaintext with incorrect size of key value", func(t *testing.T) {
		_, got := Encrypt([]byte(plaintext), []byte("hello"))
		want := "crypto/aes: invalid key size 5"

		assert.EqualError(t, got, want)
	})
}

func TestDecrypt(t *testing.T) {
	t.Run("decrypt ciphertext with 32 characters key value", func(t *testing.T) {
		got, err := Decrypt([]byte(ciphertext), []byte(key))
		want := plaintext
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, want, string(got))
	})

	t.Run("encrypt plaintext with incorrect size of key value", func(t *testing.T) {
		_, got := Decrypt([]byte(ciphertext), []byte("hello"))
		want := "crypto/aes: invalid key size 5"

		assert.EqualError(t, got, want)
	})
}
