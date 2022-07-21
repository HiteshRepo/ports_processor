package config

import "github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"

type ServerConfig struct {
	ServiceName   string       `mapstructure:"serviceName"`
	LogLevel      logger.Level `mapstructure:"logLevel"`
	PortsFilePath string       `mapstructure:"portsFilePath"`
}
