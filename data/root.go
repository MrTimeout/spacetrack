package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/model"
	l "github.com/MrTimeout/spacetrack/utils"
	"go.uber.org/zap"
)

// FileName is the template file name to write the content of space track to it
const (
	FileName   = "Spacetrack_record_"
	FileFormat = "%016d"
)

// ErrInvalidPersisterMod is returned when the persister mod passed as a parameter is incorrect
var ErrInvalidPersisterMod = errors.New("invalid persister mod")

// Persister is a type which is responsible of dumping information to one/multiple files
type Persister interface {
	Persist(string, []client.SpaceOrbitalObj) error
}

// PersisterMod represents the mod in which to persist the response from SpaceTrack.
type PersisterMod int

const (
	// OneFile is the persister that will persist the spacetrack information into a single file.
	OneFile PersisterMod = iota
	// OneFilePerRow is the persister that will persist the spacetrack information into multiple files,
	// each of them will be a row of the response.
	OneFilePerRow
)

// GetPersister returns a persister depending on the persisterMod passed as a parameter. If no
// persister is found, error is returned.
func GetPersister(pm PersisterMod, mFormat model.Format) (Persister, error) {
	var (
		p   Persister
		err error
	)

	switch pm {
	case OneFile:
		p, err = oneFile(mFormat)
	case OneFilePerRow:
		p, err = oneFilePerRow(mFormat)
	default:
		err = ErrInvalidPersisterMod
	}

	return p, err
}

// OneFilePersister is the persister that will persist the spacetrack information into a single file.
type oneFilePersister struct {
	Writer
}

func oneFile(mFormat model.Format) (Persister, error) {
	m, err := getMarshaller(mFormat)
	return oneFilePersister{WriterImpl{m}}, err
}

func (o oneFilePersister) Persist(folder string, arr []client.SpaceOrbitalObj) error {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	return o.Write(buildFilepath(time.Now().Unix(), folder), client.SpaceOrbitalObjArr{SpaceOrbitalObjs: arr})
}

// OneFilePerRowPersister
type oneFilePerRowPersister struct {
	Writer
}

func oneFilePerRow(mFormat model.Format) (Persister, error) {
	m, err := getMarshaller(mFormat)
	return oneFilePerRowPersister{WriterImpl{m}}, err
}

// Persist is used to persist all the data fetched from SpaceTrack
func (o oneFilePerRowPersister) Persist(folder string, arr []client.SpaceOrbitalObj) error {
	var (
		wg      sync.WaitGroup
		counter int32
	)

	defer func() {
		wg.Wait()
		l.Info("We have successfully persist files into system", zap.Int32("amount", counter))
	}()

	if err := cleanUp(folder); err != nil {
		return err
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}

	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := o.Write(buildFilepath(int64(i), folder), arr[i]); err != nil {
				l.Error("trying to write file to system", zap.Error(err))
			} else {
				atomic.AddInt32(&counter, 1)
			}
		}(i)
	}

	return nil
}

func buildFilepath(i int64, folder string) string {
	return filepath.Join(folder, FileName+fmt.Sprintf(FileFormat, i))
}

// TODO here we are removing all the files that match a filename ending, but we can have an html, a json...
func cleanUp(folder string) error {
	files, err := filepath.Glob(filepath.Join(folder, FileName+"*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			l.Warn("trying to delete a file while cleaning up the folder", zap.String("folder_name", folder), zap.String("file_name", f))
		}
	}

	return nil
}
