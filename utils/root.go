package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
)

//MarshalIntoFile marshal the input to a file passed as an argument
func MarshalIntoFile(input any, file string) error {
	data, err := json.Marshal(input)
	if err != nil {
		Warn("marshalling cookies", zap.Error(err))
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Warn("opening/creating the file", zap.String("file", file), zap.Error(err))
		return err
	}

	if _, err := f.Write(data); err != nil {
		Warn("writing to the file", zap.String("file", file), zap.Error(err))
		return err
	}

	return nil
}

//FileExists check if a file exists in the file system
func FileExists(filePath string) bool {
	// If we use os.O_EXCL, it always will return the ErrExist, even if it is ErrPermission.
	var err error
	if _, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666); err == nil {
		return true
	}

	// If os.ErrPermission is the error, we can't read|write to the file, so we just ignore the file
	return !errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrExist)
}

//CheckComparableIsIn allows us to check if a comparable (int, float, string) exists inside an array of comparables
func CheckComparableIsIn[K comparable](arr []K, target K) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}

//CheckStringsIsIn allows us to check if inside an array of strings, we have a "target" string, not being case sensitive
func CheckStringsIsIn(arr []string, target string) bool {
	target = strings.ToLower(target)
	for _, a := range arr {
		if strings.ToLower(a) == target {
			return true
		}
	}
	return false
}
