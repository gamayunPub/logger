package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultLevel = InfoLevel
)

type Config struct {
	Debug       bool
	Level       string
	Output      []string
	TimeEncoder string
	TimeLayout  string
}

func (c *Config) newBuilder() (*zap.Config, error) {
	var (
		config zap.Config
		err    error
	)

	if c.Debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level, err = c.getLevel()
	if err != nil {
		return nil, err
	}

	if c.TimeLayout != "" {
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(c.TimeLayout)
	}

	if len(c.Output) != 0 {
		config.OutputPaths = c.Output
	}

	return &config, nil
}

func (c *Config) getLevel() (zap.AtomicLevel, error) {
	switch c.Level {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zapcore.DebugLevel), nil
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel), nil
	case InfoLevel:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel), nil
	case WarnLevel:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel), nil
	default:
		return zap.NewAtomicLevel(), fmt.Errorf("unknown log level '%s'", c.Level)
	}
}

func (c *Config) setDefaults() {
	if c.Level == "" {
		c.Level = DefaultLevel
	}
}
