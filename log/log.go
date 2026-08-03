package log

import "go.uber.org/zap"

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal

	LogEnvDebug LogEnv = iota
	LogEnvDevelopment
	LogEnvExample
	LogEnvProduction

	_default_log_level = LevelDebug
	_default_log_env   = LogEnvDebug
)

var (
	_global_logger *MutateLogger
)

func GetGMutateLogger() *MutateLogger {
	if _global_logger == nil {
		_global_logger = NewMutateLogger(
			WithLogLevel(_default_log_level),
			WithLogEnv(_default_log_env),
		)
	}
	return _global_logger
}

func SetGMutateLogger(logger *MutateLogger) {
	_global_logger = logger
}

type Level uint8
type LogEnv uint8

type LogContext struct {
	level Level
	env   LogEnv
}

type LogContextOption func(*LogContext)

func WithLogLevel(level Level) LogContextOption {
	return func(lc *LogContext) {
		lc.level = level
	}
}

func WithLogEnv(env LogEnv) LogContextOption {
	return func(lc *LogContext) {
		lc.env = env
	}
}

type MutateLogger struct {
	lcontext   *LogContext
	zapAdapter *zapAdapter
}

func NewMutateLogger(opts ...LogContextOption) *MutateLogger {
	mLogger := &MutateLogger{
		lcontext: &LogContext{},
	}
	for _, opt := range opts {
		opt(mLogger.lcontext)
	}
	return mLogger
}

func InitMutateLogger(opts ...LogContextOption) *MutateLogger {
	mLogger := NewMutateLogger(opts...)
	_global_logger = mLogger
	return mLogger
}

func (m *MutateLogger) NewZap() (*zap.Logger, error) {
	return m.zapAdapter.New(m.lcontext)
}

func (m *MutateLogger) UseZap() (*zap.Logger, error) {
	if m.zapAdapter == nil {
		m.zapAdapter = &zapAdapter{}
	}
	return m.zapAdapter.Init(m.lcontext)
}
