package suites

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/hiteshpattanayak-tw/ports_processor/test/suites/constants"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path"
	"time"
)

const (
	PgConnectionString      = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	PgPort                  = "5432"
	PgUsername              = "postgres"
	PgPassword              = "superSecretPostgresPassword"
	PgDatabase              = "DB"
	PgContainerSetupTimeout = 10
)

type PostgresSuite struct {
	BaseSuite

	PathToScripts    string
	DBName           string
	DisableBindMount bool
}

func (p *PostgresSuite) SetupSuite() {
	p.createPgContainer()

	err := p.container.Start(p.ctx)
	p.Require().NoError(err)
}

func (p *PostgresSuite) createPgContainer() {
	getDBUrl := func(port nat.Port) string {
		username := PgUsername
		password := PgPassword
		database := p.DBName
		return fmt.Sprintf(PgConnectionString, username, password, "localhost", port.Port(), database)
	}

	pathToSeedSql := path.Join(p.rootDir, p.PathToScripts)

	request := testcontainers.ContainerRequest{
		Image:        constants.PostgresImage,
		ExposedPorts: []string{PgPort + "/tcp"},
		WaitingFor:   wait.ForSQL(PgPort+"/tcp", "postgres", getDBUrl).Timeout(PgContainerSetupTimeout * time.Second),
		Env: map[string]string{
			"POSTGRES_DB":       p.DBName,
			"POSTGRES_USER":     PgUsername,
			"POSTGRES_PASSWORD": PgPassword,
		},
	}

	if !p.DisableBindMount {
		request.Mounts = testcontainers.Mounts(testcontainers.BindMount(pathToSeedSql, "/docker-entrypoint-initdb.d"))
	}

	container, err := p.BaseSuite.GenericContainer(testcontainers.GenericContainerRequest{
		ContainerRequest: request,
	})

	p.Require().NoError(err)
	p.container = container
}

func (p *PostgresSuite) GetConnectionString(host, port string) string {
	return fmt.Sprintf(PgConnectionString, PgUsername, PgPassword, host, port, PgDatabase)
}

func (p *PostgresSuite) GetContainerMappedPort() nat.Port {
	return p.BaseSuite.GetContainerMappedPort(PgPort)
}
