package main

import (
	"os"
	"path"
	"strings"
	"testing"
)

func createFile(t *testing.T, dir, fileName string, perm os.FileMode) *os.File {
	var err error

	if strings.TrimSpace(dir) == "" {
		if dir, err = os.Getwd(); err != nil {
			t.Fatal(err)
		}
	}

	f, err := os.OpenFile(path.Join(dir, fileName), os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	return f
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
