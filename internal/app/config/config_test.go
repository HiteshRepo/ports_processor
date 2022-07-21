package config_test

import (
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/config"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type configSuite struct {
	suite.Suite
	appConfig config.AppConfig
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(configSuite))
}

func (c *configSuite) SetupTest() {
	os.Clearenv()
	configReader, err := os.Open("../../../configs/default.yaml")
	c.Require().NoError(err)
	environmentVariables := map[string]string{
		"LOG_LEVEL":              "debug",

		"DB_NAME":     "postgresdb",
		"DB_USER":     "tsdbadmin1",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "tsdbpwd",
		"DB_SCHEMA":   "public",
	}
	for environmentVariable, value := range environmentVariables {
		err := os.Setenv(environmentVariable, value)
		c.Require().NoError(err)
	}

	c.appConfig, err = config.LoadConfig(configReader)

	c.Require().NoError(err)
}

func (c *configSuite) TestServerConfig() {
	serverConfig := c.appConfig.GetServerConfig()

	c.Assert().NotNil(serverConfig, "serverConfig should not be nil")
	c.Assert().Equal("ports-processor", serverConfig.ServiceName)
	c.Assert().Equal(logger.Level("debug"), serverConfig.LogLevel)
}

func (c *configSuite) TestDatabaseConfig() {
	databaseConfig := c.appConfig.GetDatabaseConfig()

	c.Assert().Equal("postgresdb", databaseConfig.Dbname)
	c.Assert().Equal("localhost", databaseConfig.Host)
	c.Assert().Equal(5432, databaseConfig.Port)
	c.Assert().Equal("tsdbadmin1", databaseConfig.Username)
	c.Assert().Equal("tsdbpwd", databaseConfig.Password)
	c.Assert().Equal(true, databaseConfig.LogMode)
	c.Assert().Equal("disable", databaseConfig.SslMode)
	c.Assert().Equal("./scripts", databaseConfig.MigrationPath)
	c.Assert().Equal("public", databaseConfig.Schema)

	postgresConnPoolConfig := databaseConfig.Connection

	c.Assert().Equal(30, postgresConnPoolConfig.MaxOpenConnections)
	c.Assert().Equal(10, postgresConnPoolConfig.MaxIdleConnections)
	c.Assert().Equal(30, postgresConnPoolConfig.MaxIdleTime)
	c.Assert().Equal(3600, postgresConnPoolConfig.MaxLifeTime)
	c.Assert().Equal(30, postgresConnPoolConfig.TimeOut)
}
