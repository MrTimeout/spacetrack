package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
)

func MarshalIntoFile(input any, file string) error {
	data, err := json.Marshal(input)
	if err != nil {
		Logger.Warn("marshalling cookies", zap.Error(err))
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Warn("opening/creating the file", zap.String("file", file), zap.Error(err))
		return err
	}

	if _, err := f.Write(data); err != nil {
		Logger.Warn("writing to the file", zap.String("file", file), zap.Error(err))
		return err
	}

	return nil
}

func FileExists(filePath string) bool {
	// If we use os.O_EXCL, it always will return the ErrExist, even if it is ErrPermission.
	var err error
	if _, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666); err == nil {
		return true
	}

	// If os.ErrPermission is the error, we can't read|write to the file, so we just ignore the file
	return !errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrExist)
}

func CheckErr(err error) {
	if err != nil {
		Logger.Error("", zap.Error(err))
	}
}

func CheckStringIsIn(arr []string, target string) bool {
	target = strings.ToLower(target)
	for _, a := range arr {
		if strings.ToLower(a) == target {
			return true
		}
	}
	return false
}
