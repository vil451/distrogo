package data

import (
	"errors"
	"io/fs"
	"os"
)

const (
	DefaultFileMod os.FileMode = 0600
	DefaultDirMod  os.FileMode = 0744
	MainConfigFile             = "config.yml"
)

func EnsureFullPath(path string, mod os.FileMode) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(path, mod); err != nil {
			return err
		}
	}

	return nil
}
