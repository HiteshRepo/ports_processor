package config

import (
	"flag"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/config"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"io"
	"os"
	"sync"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = ""
	configFileUsage   = "this is config file path"
)

var (
	once         sync.Once
	cachedConfig AppConfig
)

type AppConfig struct {
	ServerConfig   ServerConfig      `mapstructure:"app"`
	DatabaseConfig db.DatabaseConfig `mapstructure:"db"`
}

func (a *AppConfig) GetServerConfig() ServerConfig {
	return a.ServerConfig
}

func (a *AppConfig) GetDatabaseConfig() db.DatabaseConfig {
	return a.DatabaseConfig
}

func LoadConfig(reader io.Reader) (c AppConfig, err error) {

	keysToEnvironmentVariables := map[string]string{
		"app.logLevel": "LOG_LEVEL",

		"db.name":     "DB_NAME",
		"db.user":     "DB_USER",
		"db.host":     "DB_HOST",
		"db.port":     "DB_PORT",
		"db.schema":   "DB_SCHEMA",
		"db.password": "DB_PASSWORD",
		"db.sslmode":  "DB_SSLMODE",
	}

	err = config.LoadConfig(reader, keysToEnvironmentVariables, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func ProvideAppConfig() (c AppConfig, err error) {
	once.Do(func() {
		var configFile string
		flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
		flag.Parse()

		var configReader io.ReadCloser
		configReader, err = os.Open(configFile)

		if err != nil {
			return
		}

		c, err = LoadConfig(configReader)
		if err != nil {
			return
		}

		cachedConfig = c
		_ = configReader.Close()
	})

	return cachedConfig, err
}
