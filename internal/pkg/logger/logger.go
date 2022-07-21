package logger

//go:generate mockery --name=Logger

import (
	"context"
	stdLog "log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Retrieved from build flags
var (
	Version string
	Commit  string
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(ctx context.Context) Logger
	AsStandardLogger() *stdLog.Logger
}

type ZapLogger struct {
	*zap.Logger
}

type Level string

func ProvideLogger(level Level) (l *ZapLogger, err error) {
	zl, err := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level.zapLevel()),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "@timestamp",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
			},
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		InitialFields: map[string]interface{}{
			"version": Version,
			"commit":  Commit,
		},
	}.Build()

	return &ZapLogger{zl}, err
}

func (l *ZapLogger) With(ctx context.Context) Logger {
	span, _ := tracer.SpanFromContext(ctx)
	spanCtx := span.Context()

	var fArr []zap.Field

	if spanCtx != nil {
		if traceId := spanCtx.TraceID(); traceId != 0 {
			fArr = append(fArr, zap.Uint64("dd.trace_id", spanCtx.TraceID()))
		}
		if spanId := spanCtx.SpanID(); spanId != 0 {
			fArr = append(fArr, zap.Uint64("dd.span_id", spanCtx.SpanID()))
		}
	}

	return &ZapLogger{Logger: l.Logger.With(fArr...)}
}

func (l *ZapLogger) AsStandardLogger() *stdLog.Logger {
	return zap.NewStdLog(l.Logger)
}

func (level Level) zapLevel() zapcore.Level {
	l, ok := levelMap[string(level)]
	if !ok {
		return zapcore.InfoLevel
	}

	return l
}

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}
