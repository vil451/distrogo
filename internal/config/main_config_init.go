package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func LoadMainConfig() *FileConfig {
	config := makeDefaultConfig()

	var pathToCheck []string
	for _, path := range configPaths {
		for _, ext := range configExt {
			configFile := filepath.Join(path, configName+"."+ext)
			fmt.Println(configFile)
			if _, err := os.Stat(configFile); err == nil {
				pathToCheck = append(pathToCheck, configFile)
			}
		}
	}

	err := ReadConfigFile(pathToCheck, config)
	if err != nil {
		panic(err)
	}

	return config
}

func makeDefaultConfig() *FileConfig {
	fileConfig := &FileConfig{}
	return fileConfig
}

var (
	configName  = "config"
	configExt   = []string{"yaml", "yml", "json"}
	configPaths = []string{
		"/etc/distrogo/",
		os.Getenv("XDG_CONFIG_HOME") + "/distrogo/",
		os.Getenv("HOME") + "/distrogo/",
		"./",
	}
)
