package config_test

import (
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/config"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type configSuite struct {
	suite.Suite
	tc testConfig
}

type testConfig struct {
	App struct {
		ServiceName string `mapstructure:"serviceName"`
		LogLevel    string `mapstructure:"logLevel"`
	} `mapstructure:"app"`
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(configSuite))
}

func (c *configSuite) SetupTest() {
	os.Clearenv()
	configReader, err := os.Open("test/test.yaml")
	c.Require().NoError(err)
	binding := map[string]string{
		"app.logLevel": "LOG_LEVEL",
	}

	err = os.Setenv("LOG_LEVEL", "info")
	c.Require().NoError(err)

	err = config.LoadConfig(configReader, binding, &c.tc)
	c.Require().NoError(err)
}

func (c *configSuite) TestLoadConfig() {
	c.Assert().Equal("port-processor", c.tc.App.ServiceName)
	c.Assert().Equal("info", c.tc.App.LogLevel)
}
