package zapshit

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Main() {
	demo2()
}

type M struct {
	Name string
	Age  uint8
}

func demo1() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	m := NewM("MJ", 23)
	if m == nil {
		return
	}
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		// zap.String("url", m.Name),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
		zap.Any("m", *m),
	)
}

func NewM(name string, age uint8) *M {
	return &M{
		Name: name,
		Age:  age,
	}
}

//go:generate stringer -type=SearcherName
type SearcherName int

const (
	Err SearcherName = iota
	InterceptWord
	Greeting
	TabledKnowledge
	Knowledge
	ChatGPT
	Helplook
)

func demo2() {
	logger := createLogger()
	defer logger.Sync()

	s := TabledKnowledge

	logger.Info("Hello from Zap!", zap.Any("search_level", s), zap.Duration("hello", 60987654321))
}

func createLogger() *zap.Logger {
	// encoderCfg := zap.NewProductionEncoderConfig()
	// encoderCfg.TimeKey = "T"
	// encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     newEncoderCfg(),
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func newEncoderCfg() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		// EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		// EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}
