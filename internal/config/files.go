package config

import (
	"distrogo/internal/config/data"
	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	distrgoConfigDir = "CONFIG_DIR"
	AppName          = "distrogo"
	LogsFile         = "distrogo.log"
)

var (
	AppConfigDir   string
	AppDumpsDir    string
	AppConfigFile  string
	AppContextsDir string
)

func initXDGLocs() error {
	var err error

	AppConfigDir, err = xdg.ConfigFile(AppName)
	if err != nil {
		return err
	}

	AppConfigFile, err = xdg.ConfigFile(filepath.Join(AppName, data.MainConfigFile))
	if err != nil {
		return err
	}

	dataDir, err := xdg.DataFile(AppName)
	if err != nil {
		return err
	}
	AppContextsDir = filepath.Join(dataDir, "distros")
	if err := data.EnsureFullPath(AppContextsDir, data.DefaultDirMod); err != nil {
		log.Warn().Err(err).Msgf("No context dir")
	}

	return nil
}

func initEnvLoc() error {
	AppConfigDir = os.Getenv(distrgoConfigDir)
	if err := data.EnsureFullPath(AppConfigDir, data.DefaultDirMod); err != nil {
		return err
	}
	return nil
}

func isEnvSet(env string) bool {
	return os.Getenv(env) != ""
}

func InitLogLocs() error {
	var _ string
	tmpDir, err := UserTmpDir()
	if err != nil {
		return err
	}
	_ = tmpDir
	if err := data.EnsureFullPath(tmpDir, data.DefaultDirMod); err != nil {
		return err
	}
	AppLogFile = filepath.Join(AppLogFile, LogsFile)
	return nil
}

func UserTmpDir() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(os.TempDir(), current.Username, AppName)
	return dir, nil
}
