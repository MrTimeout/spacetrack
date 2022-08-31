package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/model"
	l "github.com/MrTimeout/spacetrack/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	l.Configure(l.NewLogger(true, l.DebugLevel))
	os.Exit(m.Run())
}

func createTmpFolder(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path, err := os.MkdirTemp(dir, "oneFilePersister")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(path) //golint:errcheck
	})

	return path
}

func TestGetPersister(t *testing.T) {
	for _, each := range []struct {
		description string
		input       PersisterMod
		want        Persister
		wantErr     error
	}{
		{
			description: "persisterMod of one file returns the one file persister",
			input:       OneFile,
			want:        oneFilePersister{WriterImpl{XMLMarshaller{}}},
		},
		{
			description: "persisterMod of one file per rowreturns the one file per row persister",
			input:       OneFilePerRow,
			want:        oneFilePerRowPersister{WriterImpl{XMLMarshaller{}}},
		},
		{
			description: "persisterMod of one file returns the one file persister",
			input:       3,
			want:        nil,
			wantErr:     ErrInvalidPersisterMod,
		},
	} {

		t.Run(each.description, func(t *testing.T) {
			got, gotErr := GetPersister(each.input, model.Xml)

			assert.Equal(t, each.want, got)
			assert.ErrorIs(t, each.wantErr, gotErr)
		})
	}
}

func TestOneFilePersister(t *testing.T) {
	for _, each := range []struct {
		description string
		pm          PersisterMod
		arr         []client.SpaceOrbitalObj
		want        int
	}{
		{
			description: "one file persister",
			pm:          OneFile,
			arr: []client.SpaceOrbitalObj{
				{Comment: "here is a comment", Originator: "originator1"},
				{Comment: "here is another comment", Originator: "originator2"},
			},
			want: 1,
		},
		{
			description: "one file per row persister",
			pm:          OneFilePerRow,
			arr: []client.SpaceOrbitalObj{
				{Comment: "here is a comment", Originator: "originator1"},
				{Comment: "here is another comment", Originator: "originator2"},
				{Comment: "here is another comment", Originator: "originator3"},
				{Comment: "here is another comment", Originator: "originator4"},
				{Comment: "here is another comment", Originator: "originator5"},
			},
			want: 5,
		},
	} {
		t.Run(each.description, func(t *testing.T) {
			var inputFolder = createTmpFolder(t)

			p, err := GetPersister(each.pm, model.Xml)
			if err != nil {
				t.Fatal(err)
			}

			err = p.Persist(inputFolder, each.arr)
			if err != nil {
				t.Fatal(err)
			}

			files, err := filepath.Glob(filepath.Join(inputFolder, FileName) + "*")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, each.want, len(files))
		})
	}
}
