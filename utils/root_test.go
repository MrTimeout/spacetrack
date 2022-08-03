package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrMarshalling = errors.New("utils/root_test: DumMarshaller can't marshal")

//DumbMarshaller interface Marshaller
type DumbMarshaller struct {
	Value string `json:"value"`
}

func (d DumbMarshaller) MarshalJSON() ([]byte, error) {
	return nil, ErrMarshalling
}

func createTmpWithPerm(t *testing.T, mod os.FileMode) string {
	f, err := os.CreateTemp("./", "tmpWithNoperm")
	if err != nil {
		t.Fatal(err)
	}

	// Changing permissions to only read
	if err = f.Chmod(mod); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})
	return f.Name()
}

func TestMarshalIntoFile(t *testing.T) {
	t.Run("marshal a nil value", func(t *testing.T) {
		got := MarshalIntoFile(DumbMarshaller{Value: "some value"}, "")
		want := ErrMarshalling.Error()

		assert.Contains(t, got.Error(), want)
	})

	t.Run("trying to open a file without right permission", func(t *testing.T) {
		got := MarshalIntoFile(struct{ Value string }{Value: "amazing"}, createTmpWithPerm(t, 0444))

		assert.ErrorIs(t, got, os.ErrPermission)
	})

	t.Run("marshal right value", func(t *testing.T) {
		f := createTmpWithPerm(t, 0644)
		err := MarshalIntoFile(struct{ Value string }{Value: "amazing"}, f)
		if err != nil {
			t.Fatal(err)
		}

		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "{\"Value\":\"amazing\"}", string(b))
	})
}

func TestFileExists(t *testing.T) {
	t.Run("file doesn't exist when err permission", func(t *testing.T) {
		assert.False(t, FileExists(createTmpWithPerm(t, 0444)))
	})

	t.Run("file already exists because it really exists", func(t *testing.T) {
		assert.True(t, FileExists(createTmpWithPerm(t, 0644)))
	})

	t.Run("file doesn't exist and we can create it so the result is true", func(t *testing.T) {
		f := "./randomfile"
		t.Cleanup(func() {
			os.Remove(f)
		})

		assert.True(t, FileExists(f))
	})
}
