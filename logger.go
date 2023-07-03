package logger

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogKey string

const (
	DebugLevel = "debug"
	ErrorLevel = "error"
	InfoLevel  = "info"
	WarnLevel  = "warn"
)

const (
	TimeEncoderEpoch   = "epoch"
	TimeEncoderISO8601 = "ISO8601"
)

func New(cfg *Config, appName, appVersion string) (*Logger, error) {
	cfg.setDefaults()

	builder, err := cfg.newBuilder()
	if err != nil {
		return nil, err
	}

	logger, err := builder.Build()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	name := zap.String("app_name", appName)
	version := zap.String("app_version", appVersion)
	level := zap.AddCallerSkip(1)

	return wrap(logger.With(name, version).WithOptions(level)), nil
}

func wrap(logger *zap.Logger) *Logger {
	return &Logger{logger: logger.Sugar()}
}

type Logger struct {
	logger    *zap.SugaredLogger
	ctxFields []LogKey
}

func (log *Logger) SetCtxFields(fields ...LogKey) {
	log.ctxFields = fields
}

func (log *Logger) AddCallerSkip(n int) {
	log.logger = log.logger.Desugar().WithOptions(zap.AddCallerSkip(n)).Sugar()
}

func (log *Logger) AddHooks(hooks ...func(zapcore.Entry) error) {
	log.logger = log.logger.Desugar().WithOptions(zap.Hooks(hooks...)).Sugar()
}

func (log *Logger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Debugw(msg, log.addCtxFields(ctx, keysAndValues...)...)
}

func (log *Logger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Infow(msg, log.addCtxFields(ctx, keysAndValues...)...)
}

func (log *Logger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Warnw(msg, log.addCtxFields(ctx, keysAndValues...)...)
}

func (log *Logger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Errorw(msg, log.addCtxFields(ctx, keysAndValues...)...)
}

func (log *Logger) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Fatalw(msg, log.addCtxFields(ctx, keysAndValues...)...)
}

func (log *Logger) addCtxFields(ctx context.Context, keysAndValues ...interface{}) []interface{} {
	result := make([]interface{}, 0)

	for _, field := range log.ctxFields {
		if v := ctx.Value(field); v != nil {
			result = append(result, string(field), v)
		}
	}

	return append(result, keysAndValues...)
}
