package migrations_test

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db/migrations"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"github.com/hiteshpattanayak-tw/ports_processor/test/suites"
	"os"
	"path"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	connectionString = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	username         = "postgres"
	password         = "superSecretPostgresPassword"
)

type migrationsSuite struct {
	suites.FullSuite

	database string
	password string

	postgresDB *gorm.DB

	migrator *migrations.Migrator
}

func TestMigrationSuite(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	cnf := suites.FullSuiteConfig{
		RootDir:        path.Join(workingDirectory, "../../../../"),
		EnablePostgres: true,
	}

	suite.Run(t, &migrationsSuite{
		FullSuite: suites.NewFullSuite(cnf),
	})
}

func (suite *migrationsSuite) SetupTest() {
	postgresSuite := suite.Suites[suites.PostgresSuiteId]
	postgresHost := postgresSuite.GetContainerHost()
	p := postgresSuite.GetContainerMappedPort()

	postgresPort, err := strconv.Atoi(p.Port())
	suite.Require().NoError(err)

	postgresDBConfig := db.DatabaseConfig{
		Dbname:        suite.database,
		Username:      "postgres",
		Password:      suite.password,
		Host:          "localhost",
		Port:          postgresPort,
		LogMode:       true,
		SslMode:       "disable",
		MigrationPath: "testdata/migrations",
	}

	l, err := logger.ProvideLogger("")
	suite.Require().NoError(err)

	dsn := fmt.Sprintf(connectionString, username, password, postgresHost, p.Port(), suite.database)
	suite.postgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	suite.migrator, err = migrations.ProvideMigrator(postgresDBConfig, suite.postgresDB, l)
	suite.Require().NoError(err)
}

func (suite *migrationsSuite) TestAppliesUpPostgresMigrations() {
	suite.Require().False(suite.migrationsApplied(suite.postgresDB))

	suite.migrator.RunMigrations()

	suite.Require().True(suite.migrationsApplied(suite.postgresDB))
}

func (suite *migrationsSuite) migrationsApplied(db *gorm.DB) bool {
	var tableExists bool
	err := db.Raw(`SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'first_table');`).Row().Scan(&tableExists)
	suite.Require().NoError(err)

	var indexExists bool
	err = db.Raw(`SELECT EXISTS (SELECT FROM pg_class WHERE relname = 'name_idx');`).Row().Scan(&indexExists)
	suite.Require().NoError(err)

	return tableExists && indexExists
}
