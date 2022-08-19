package data

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/utils"
	"go.uber.org/zap"
)

// FileName is the template file name to write the content of space track to it
const FileName = "Spacetrack_record_"

var wg sync.WaitGroup
var counter int32

// Persist is used to persist all the data fetched from SpaceTrack
func Persist(folder string, spaceOrbitalObjectArr []client.SpaceOrbitalObj) error {
	defer func() {
		wg.Wait()
		utils.Info("We have successfully persist files into system", zap.Int32("amount", counter))
	}()

	if err := cleanUp(folder); err != nil {
		return err
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}

	for i, spaceOrbitalObj := range spaceOrbitalObjectArr {
		wg.Add(1)
		go write(i, folder, spaceOrbitalObj)
	}

	return nil
}

func write(i int, folder string, row client.SpaceOrbitalObj) {
	defer wg.Done()
	file := filepath.Join(folder, FileName+fmt.Sprintf("%08d.xml", i))
	buff, err := xml.Marshal(row)
	if err != nil {
		utils.Warn("trying to serializing xml", zap.String("file", file), zap.Error(err))
		return
	}

	if err := os.WriteFile(file, buff, 0644); err != nil {
		utils.Warn("writing xml to file", zap.String("file", file), zap.Error(err))
		return
	}

	atomic.AddInt32(&counter, 1)
}

func cleanUp(folder string) error {
	files, err := filepath.Glob(filepath.Join(folder, FileName+"*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			utils.Warn("trying to delete a file while cleaning up the folder", zap.String("folder", folder), zap.String("file", f))
		}
	}

	return nil
}
