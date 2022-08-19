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
	// ErrCipherTextTooShort is thrown when cipher text is too short
	ErrCipherTextTooShort = errors.New("utils/credentials: ciphertext too short")

	// ... random bytes
	iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
)

// Encrypt will encrypt the plain text using the passphrase passed as a parameter
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

// Encode will encode a buffer with base64
func Encode(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

// Decrypt will decrypt the cipher text using the passphrase passed as a parameter
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

// Decode will decode a buffer with base64
func Decode(b []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}

// ReadOnlyFilePassphrase will read the passphrase from a file
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

	b, err = Decode(b)
	if err != nil {
		return "", err
	}

	if len(b) != 32 {
		return "", ErrCipherTextTooShort
	}

	return string(b), nil
}

// WritePassphraseToFile will write the passphrase to the file
func WritePassphraseToFile(filename string, passphrase []byte) error {
	var (
		f   *os.File
		err error
	)

	if f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664); err != nil {
		return err
	}

	_, err = f.Write(passphrase)

	return err
}
