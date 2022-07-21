package logger_test

import (
	"context"
	"encoding/json"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"github.com/hiteshpattanayak-tw/ports_processor/test"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	stdLog "log"
	"strings"
	"testing"
	"time"
)

type loggerSuite struct {
	suite.Suite
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(loggerSuite))
}

func (suite *loggerSuite) SetupTest() {
	logger.Version = "1.2.3"
	logger.Commit = "b2c3e317682929a3255a5e6433ccc5be"
}

func (suite *loggerSuite) TestShouldLogVersionNumberAndCommit() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("info")
		suite.Require().NoError(err)

		logger.Debug("This should not be logged, since the level is info")
		logger.Info("Test")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("INFO", parsedLog["level"])
	suite.Assert().Equal("Test", parsedLog["message"])
	suite.Assert().Equal("1.2.3", parsedLog["version"])
	suite.Assert().Equal("b2c3e317682929a3255a5e6433ccc5be", parsedLog["commit"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestDefaultLogLevelShouldBeInfo() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("")
		suite.Require().NoError(err)
		logger.Debug("This should not be logged, since the level is info")
		logger.Info("Test")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("INFO", parsedLog["level"])
	suite.Assert().Equal("Test", parsedLog["message"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestDebugLogging() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("debug")
		suite.Require().NoError(err)
		logger.Debug("Debug message")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("DEBUG", parsedLog["level"])
	suite.Assert().Equal("Debug message", parsedLog["message"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestInfoLogging() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("debug")
		suite.Require().NoError(err)
		logger.Info("Info message")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("INFO", parsedLog["level"])
	suite.Assert().Equal("Info message", parsedLog["message"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestWarnLogging() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("")
		suite.Require().NoError(err)
		logger.Warn("Warn message")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("WARN", parsedLog["level"])
	suite.Assert().Equal("Warn message", parsedLog["message"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestErrorLogging() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("debug")
		suite.Require().NoError(err)
		logger.Error("Error message")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("ERROR", parsedLog["level"])
	suite.Assert().Equal("Error message", parsedLog["message"])
	suite.Assert().NotEmpty(parsedLog["caller"])
}

func (suite *loggerSuite) TestLogsWithUTCTime() {
	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("debug")
		suite.Require().NoError(err)
		logger.Info("Test")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	rawTimestamp := parsedLog["@timestamp"].(string)
	parsedTimestamp, err := time.Parse("2006-01-02T15:04:05Z0700", rawTimestamp)
	suite.Require().NoError(err)

	suite.Assert().True(strings.HasSuffix(rawTimestamp, "Z"))
	suite.Assert().Equal(time.UTC, parsedTimestamp.Location())
}

func (suite *loggerSuite) TestLoggingWithContextForATraceAndSpan() {
	tracer.Start()

	_, ctx := tracer.StartSpanFromContext(context.Background(),
		"Some Span Operation",
		tracer.WithSpanID(123),
	)

	rawLog := test.CaptureStandardOut(suite.T(), func() {
		logger, err := logger.ProvideLogger("debug")
		suite.Require().NoError(err)
		logger.With(ctx).Info("Test")
	})

	parsedLog := parseLog(suite.T(), rawLog)
	suite.Assert().Equal("INFO", parsedLog["level"])
	suite.Assert().Equal("Test", parsedLog["message"])
	suite.Assert().Equal(float64(123), parsedLog["dd.span_id"])
	suite.Assert().Equal(float64(123), parsedLog["dd.trace_id"])
	suite.Assert().NotEmpty(parsedLog["caller"])

	tracer.Stop()
}

func (suite *loggerSuite) TestCanConvertToTheStandardLibraryLogger() {
	logger, err := logger.ProvideLogger("info")
	suite.Require().NoError(err)

	suite.Require().NotNil(logger.AsStandardLogger())
	suite.Require().IsType(&stdLog.Logger{}, logger.AsStandardLogger())
}

func parseLog(t *testing.T, log string) map[string]interface{} {
	var parsedLog map[string]interface{}
	err := json.Unmarshal([]byte(log), &parsedLog)
	require.NoError(t, err)
	return parsedLog
}
