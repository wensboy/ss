package log

import (
	"go.uber.org/zap"
)

type zapAdapter struct {
	zapLogger *zap.Logger
}

func (z *zapAdapter) New(lc *LogContext) (*zap.Logger, error) {
	var err error
	z.zapLogger, err = zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	switch lc.env {
	case LogEnvProduction:
		z.zapLogger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case LogEnvExample:
		z.zapLogger = zap.NewExample()
	}
	return z.zapLogger, nil
}

func (z *zapAdapter) Init(lc *LogContext) (*zap.Logger, error) {
	if z.zapLogger != nil {
		return z.zapLogger, nil
	}
	return z.New(lc)
}
