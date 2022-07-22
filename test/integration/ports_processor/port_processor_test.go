package ports_processor_test

import (
	"context"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/config"
	dbModel "github.com/hiteshpattanayak-tw/ports_processor/internal/app/db/migrations/model"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db/migrations"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"github.com/hiteshpattanayak-tw/ports_processor/test/suites"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"path"
	"strconv"
	"testing"
)

type portsProcessorSuite struct {
	suites.FullSuite
	ctx    context.Context
	cancel context.CancelFunc

	app  app.App
	db   *gorm.DB
}

func TestPortsProcessorSuite(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	cnf := suites.FullSuiteConfig{
		RootDir:        path.Join(workingDirectory, "../../../"),
		EnablePostgres: true,
	}

	suite.Run(t, &portsProcessorSuite{
		FullSuite: suites.NewFullSuite(cnf),
	})
}

func (s *portsProcessorSuite) SetupSuite() {
	s.FullSuite.SetupSuite()

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.initPortsProcessor()
	s.servePortsProcessor()
}

func (s *portsProcessorSuite) TestPortsProcessor() {
	var ports []dbModel.Port
	result := s.db.Find(&ports)
	s.Require().Nil(result.Error)
	s.Assert().True(len(ports) == 2)
}

func (s *portsProcessorSuite) initPortsProcessor() {
	workingDirectory, _ := os.Getwd()
	pathToConfig := path.Join(workingDirectory, "../../../configs/local.yaml")

	configReader, err := os.Open(pathToConfig)
	s.Require().NoError(err)

	appConfig, err := config.LoadConfig(configReader)
	s.Require().NoError(err)
	appConfig.ServerConfig.PortsFilePath = "test/integration/ports_processor/ports.json"

	l, err := logger.ProvideLogger("")
	s.Require().NoError(err)

	postgresSuite := s.Suites[suites.PostgresSuiteId]
	host := postgresSuite.GetContainerHost()
	p := postgresSuite.GetContainerMappedPort()
	dbPort, err := strconv.Atoi(p.Port())
	s.Require().NoError(err)

	dbCnf := appConfig.GetDatabaseConfig()
	dbCnf.Host = host
	dbCnf.Port = dbPort
	dbCnf.Password = "superSecretPostgresPassword"
	dbCnf.Dbname = "test_db"
	dbCnf.MigrationPath = "../../../test/scripts"

	db, err := db.ProvideDatabase(dbCnf, "ports-processor")
	s.Require().NoError(err)
	s.db = db

	migrator, err := migrations.ProvideMigrator(dbCnf, db, l)
	s.Require().NoError(err)

	portRepository := repository.ProvidePortRepository(db)

	stream := json_processor.ProvideJSONStream()

	s.app = app.App{
		Ctx:        s.ctx,
		Cancel:     s.cancel,
		Migrator:   migrator,
		Logger:     l,
		PortRepo:   portRepository,
		AppConfig:  appConfig,
		JsonStream: stream,
	}
}

func (s *portsProcessorSuite) servePortsProcessor() {
	workingDirectory, err := os.Getwd()
	s.Require().NoError(err)

	rootDir := path.Join(workingDirectory, "../../../")

	s.app.Start(func(err error) {
		s.Require().NoError(err)
	}, rootDir)
}
