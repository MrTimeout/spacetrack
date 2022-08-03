package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var (
	ErrCipherTextTooShort = errors.New("utils/credentials: ciphertext too short")

	iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
)

func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Ciphertext FeedBack
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return Encode(ciphertext), nil
}

func Encode(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext, err = Decode(ciphertext)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

func Decode(b []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}

func ReadOnlyFilePassphrase(filename string) (string, error) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.OpenFile(filename, os.O_RDONLY, 0400); err != nil {
		return "", err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	if len(b) != 32 {
		return "", ErrCipherTextTooShort
	}

	return string(b), nil
}
