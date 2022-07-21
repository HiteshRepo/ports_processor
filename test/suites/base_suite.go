package suites

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/hiteshpattanayak-tw/ports_processor/test/suites/constants"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type BaseSuite struct {
	suite.Suite
	ctx       context.Context
	network   *testcontainers.DockerNetwork
	container testcontainers.Container
	rootDir   string
}

func (s *BaseSuite) GetContainerHost() string {
	host, err := s.container.Host(s.ctx)
	s.Require().NoError(err)

	return host
}

func (s *BaseSuite) GetContainerMappedPort(portStr string) nat.Port {
	port, err := s.container.MappedPort(s.ctx, nat.Port(portStr))
	s.Require().NoError(err)

	return port
}

func (s *BaseSuite) SetCtx(ctx context.Context) {
	s.ctx = ctx
}

func (s *BaseSuite) SetNetwork(network *testcontainers.DockerNetwork) {
	s.network = network
}

func (s *BaseSuite) SetRootDir(rootDir string) {
	s.rootDir = rootDir
}

func (s *BaseSuite) TearDownSuite() {
	s.TerminateContainer(s.container)
}

func (s *BaseSuite) TerminateContainer(container testcontainers.Container) {
	if container != nil {
		if err := container.Terminate(s.ctx); err != nil {
			s.Require().NoError(err)
		}
	}
}

func (s *BaseSuite) GenericContainer(request testcontainers.GenericContainerRequest) (testcontainers.Container, error) {
	request.ReaperImage = constants.ReaperImage
	return testcontainers.GenericContainer(s.ctx, request)
}
