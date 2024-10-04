package config

import "sync"

/*
 * If u need to add some new fields to config then firstly u need to add it co Config struct like non exported and
 * then like exported. Then u need to add getter for field in FileConfig. And at last modify mapConfig method.
 */

type FileConfig struct {
	Common struct {
		LogLevel int `json:"logLevel" yaml:"logLevel" env:"LOG_LEVEL"`
	} `json:"common" yaml:"common"`
	Docker struct {
		Host      string `json:"docker-host" yaml:"docker-host" env:"DOCKER_HOST"`
		CertPath  string `json:"certPath" yaml:"certPath" env:"DOCKER_CERT_PATH"`
		TLSVerify bool   `json:"tlsVerify" yaml:"tlsVerify" env:"DOCKER_TLS_VERIFY"`
		Registry  string `json:"registry" yaml:"registry" env:"DOCKER_REGISTRY"`
	} `json:"docker" yaml:"docker"`
	Volumes     []string `json:"volumes" yaml:"volumes"`
	EnvContexts []string `json:"envContext" yaml:"envContext"`
}

type Config struct {
	common      common
	docker      docker
	volumes     []string
	envContexts []string
}

type common struct {
	logLevel int
}

type docker struct {
	host      string
	certPath  string
	tlsVerify bool
	registry  string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

func loadConfig() *Config {
	fileConfig := LoadMainConfig()
	return mapConfig(fileConfig)
}

// Reflection not working for non exported fields, so we need to map it manually. Every time we add new Field to structs we need modify this method too.
func mapConfig(fileConfig *FileConfig) *Config {
	return &Config{
		common: common{
			logLevel: fileConfig.Common.LogLevel,
		},
		docker: docker{
			host:      fileConfig.Docker.Host,
			certPath:  fileConfig.Docker.CertPath,
			tlsVerify: fileConfig.Docker.TLSVerify,
			registry:  fileConfig.Docker.Registry,
		},
		volumes:     fileConfig.Volumes,
		envContexts: fileConfig.EnvContexts,
	}
}

func (c *Config) GetLogLevel() int {
	return c.common.logLevel
}

func (c *Config) GetDockerHost() string {
	return c.docker.host
}

func (c *Config) GetDockerCertPath() string {
	return c.docker.certPath
}

func (c *Config) GetDockerTLSVerify() bool {
	return c.docker.tlsVerify
}

func (c *Config) GetDockerRegistry() string {
	return c.docker.registry
}

func (c *Config) GetVolumes() []string {
	return c.volumes
}

func (c *Config) GetEnvContexts() []string {
	return c.envContexts
}
