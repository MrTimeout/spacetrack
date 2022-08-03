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

func createTmpWithPerm(t *testing.T, mod os.FileMode) *os.File {
	f, err := os.CreateTemp("./", "tmpFile")
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
	return f
}

func TestMarshalIntoFile(t *testing.T) {
	t.Run("marshal a nil value", func(t *testing.T) {
		got := MarshalIntoFile(DumbMarshaller{Value: "some value"}, "")
		want := ErrMarshalling.Error()

		assert.Contains(t, got.Error(), want)
	})

	t.Run("trying to open a file without right permission", func(t *testing.T) {
		got := MarshalIntoFile(struct{ Value string }{Value: "amazing"}, createTmpWithPerm(t, 0444).Name())

		assert.ErrorIs(t, got, os.ErrPermission)
	})

	t.Run("marshal right value", func(t *testing.T) {
		f := createTmpWithPerm(t, 0644)
		err := MarshalIntoFile(struct{ Value string }{Value: "amazing"}, f.Name())
		if err != nil {
			t.Fatal(err)
		}

		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "{\"Value\":\"amazing\"}", string(b))
	})
}

func TestFileExists(t *testing.T) {
	t.Run("file doesn't exist when err permission", func(t *testing.T) {
		assert.False(t, FileExists(createTmpWithPerm(t, 0444).Name()))
	})

	t.Run("file already exists because it really exists", func(t *testing.T) {
		assert.True(t, FileExists(createTmpWithPerm(t, 0644).Name()))
	})

	t.Run("file doesn't exist and we can create it so the result is true", func(t *testing.T) {
		f := "./randomfile"
		t.Cleanup(func() {
			os.Remove(f)
		})

		assert.True(t, FileExists(f))
	})
}

func TestCheckComparableIsIn(t *testing.T) {
	t.Run("check if comparable int is inside array being true", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}

		for _, i := range input {
			got := CheckComparableIsIn(input, i)

			assert.True(t, got)
		}
	})

	t.Run("check if comparable string is inside array being false", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		got := CheckComparableIsIn(input, "d")

		assert.False(t, got)
	})

	t.Run("check if comparable string is inside empty array being false", func(t *testing.T) {
		input := []string{}
		got := CheckComparableIsIn(input, "d")

		assert.False(t, got)
	})
}

func TestCheckStringsIsIn(t *testing.T) {
	t.Run("check if string lower case is inside string array being true", func(t *testing.T) {
		input := []string{"A", "B", "C"}
		got := CheckStringsIsIn(input, "b")

		assert.True(t, got)
	})

	t.Run("check if string upper case is inside string array being true", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		got := CheckStringsIsIn(input, "B")

		assert.True(t, got)
	})

	t.Run("check if string is inside empty string", func(t *testing.T) {
		input := []string{}
		got := CheckStringsIsIn(input, "A")

		assert.False(t, got)
	})

	t.Run("check if string is inside string being false", func(t *testing.T) {
		input := []string{"a", "B", "C"}
		got := CheckStringsIsIn(input, "D")

		assert.False(t, got)
	})
}
