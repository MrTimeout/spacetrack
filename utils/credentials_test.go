package utils

import (
	"strings"
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

func TestReadOnlyFilePassphrase(t *testing.T) {
	t.Run("read only file passphrase when file doesn't exist", func(t *testing.T) {
		_, got := ReadOnlyFilePassphrase("./notexistentfile")

		assert.ErrorContains(t, got, "no such file or directory")
	})

	t.Run("read only file passphrase with passphrase length different than 32", func(t *testing.T) {
		f := createTmpWithPerm(t, 0666)
		if _, err := f.Write(Encode([]byte("random data"))); err != nil {
			t.Fatal(err)
		}

		_, got := ReadOnlyFilePassphrase(f.Name())

		assert.ErrorIs(t, ErrCipherTextTooShort, got)
	})

	t.Run("read only file passphrase correctly", func(t *testing.T) {
		f := createTmpWithPerm(t, 0666)
		if _, err := f.Write(Encode([]byte(strings.Repeat("a", 32)))); err != nil {
			t.Fatal(err)
		}

		passphrase, err := ReadOnlyFilePassphrase(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, strings.Repeat("a", 32), passphrase)
	})
}
