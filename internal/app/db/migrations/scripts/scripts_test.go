package scripts_test

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/ports_processor/test/suites"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path"
	"testing"
)

const (
	connectionString = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	username         = "postgres"
	password         = "superSecretPostgresPassword"
	database         = "test_db"
)

type migrationsSuite struct {
	suites.FullSuite
	db *gorm.DB
}

type rowInfo struct {
	ColumnName string
	DataType   string
}

func TestMigrationSuite(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	cnf := suites.FullSuiteConfig{
		RootDir:         path.Join(workingDirectory, "../../../../../"),
		PathToDbScripts: "internal/app/db/migrations/scripts",
		EnablePostgres:  true,
	}

	suite.Run(t, &migrationsSuite{
		FullSuite: suites.NewFullSuite(cnf),
	})
}

func (suite *migrationsSuite) SetupTest() {
	var err error

	postgresSuite := suite.Suites[suites.PostgresSuiteId]
	host := postgresSuite.GetContainerHost()
	p := postgresSuite.GetContainerMappedPort()

	dsn := fmt.Sprintf(connectionString, username, password, host, p.Port(), database)
	suite.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)
}

func (suite *migrationsSuite) TearDownTest() {
	suite.db.Raw(`DELETE FROM ports;`).Scan(nil)
}

func (suite *migrationsSuite) TestShouldCreatePortsTable() {
	var tableInfo []rowInfo

	suite.db.Raw(`select column_name , data_type FROM information_schema."columns" c 
	WHERE table_name = 'ports' order by column_name;`).Scan(&tableInfo)

	suite.Require().Equal(10, len(tableInfo))

	suite.Assert().Equal("alias", tableInfo[0].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[0].DataType)

	suite.Assert().Equal("city", tableInfo[1].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[1].DataType)

	suite.Assert().Equal("code", tableInfo[2].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[2].DataType)

	suite.Assert().Equal("coordinates", tableInfo[3].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[3].DataType)

	suite.Assert().Equal("country", tableInfo[4].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[4].DataType)

	suite.Assert().Equal("name", tableInfo[5].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[5].DataType)

	suite.Assert().Equal("province", tableInfo[6].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[6].DataType)

	suite.Assert().Equal("regions", tableInfo[7].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[7].DataType)

	suite.Assert().Equal("timezone", tableInfo[8].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[8].DataType)

	suite.Assert().Equal("unlocs", tableInfo[9].ColumnName)
	suite.Assert().Equal("character varying", tableInfo[9].DataType)
}
