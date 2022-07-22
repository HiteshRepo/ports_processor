package repository_test

import (
	"fmt"
	dbModel "github.com/hiteshpattanayak-tw/ports_processor/internal/app/db/migrations/model"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
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

type portsRepositorySuite struct {
	suites.FullSuite
	db       *gorm.DB
	portRepo repository.PortRepository
}

func TestPortRepositorySuite(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	cnf := suites.FullSuiteConfig{
		RootDir:         path.Join(workingDirectory, "../../../"),
		PathToDbScripts: "internal/app/db/migrations/scripts",
		EnablePostgres:  true,
	}

	suite.Run(t, &portsRepositorySuite{
		FullSuite: suites.NewFullSuite(cnf),
	})
}

func (suite *portsRepositorySuite) SetupTest() {
	var err error

	postgresSuite := suite.Suites[suites.PostgresSuiteId]
	host := postgresSuite.GetContainerHost()
	p := postgresSuite.GetContainerMappedPort()

	dsn := fmt.Sprintf(connectionString, username, password, host, p.Port(), database)
	suite.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	suite.portRepo = repository.ProvidePortRepository(suite.db)
}

func (suite *portsRepositorySuite) TearDownTest() {
	suite.db.Raw(`DELETE FROM ports;`).Scan(nil)
}

func (suite *portsRepositorySuite) TestUpsertPortShouldInsertPort() {
	port := &dbModel.Port{
		Name:        "TestName",
		City:        "TestCity",
		Country:     "TestCountry",
		Alias:       "TestAlias1,TestAlias2",
		Regions:     "TestReg1,TestReg2",
		Coordinates: "TestCoord1,TestCoord2",
		Province:    "TestProvince",
		Timezone:    "TestTZ",
		Unlocs:      "TestUnlocs1,TestUnlocs2",
		Code:        "0000111",
	}
	err := suite.portRepo.UpsertPort(suite.Ctx, "ports", port)
	suite.Require().NoError(err)

	var dbPort dbModel.Port
	result := suite.db.Where("name = ?", "TestName").Last(&dbPort)
	suite.Require().Nil(result.Error)
	suite.Assert().Equal(*port, dbPort)
}

func (suite *portsRepositorySuite) TestUpsertPortShouldUpdatePort() {
	port := &dbModel.Port{
		Name:        "TestName",
		City:        "TestCity",
		Country:     "TestCountry",
		Alias:       "TestAlias1,TestAlias2",
		Regions:     "TestReg1,TestReg2",
		Coordinates: "TestCoord1,TestCoord2",
		Province:    "TestProvince",
		Timezone:    "TestTZ",
		Unlocs:      "TestUnlocs1,TestUnlocs2",
		Code:        "0000111",
	}
	err := suite.portRepo.UpsertPort(suite.Ctx, "ports", port)
	suite.Require().NoError(err)

	port.City = "TestCity2"

	err = suite.portRepo.UpsertPort(suite.Ctx, "ports", port)
	suite.Require().NoError(err)

	var dbPort dbModel.Port
	result := suite.db.Where("name = ?", "TestName").Last(&dbPort)
	suite.Require().Nil(result.Error)
	suite.Assert().True(dbPort.City == "TestCity2")
}
