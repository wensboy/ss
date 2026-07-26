package log

import "testing"

func Test_MutateLogger(t *testing.T) {
	_ = InitMutateLogger(WithLogLevel(LevelDebug), WithLogEnv(LogEnvExample))
	logger, err := _global_logger.UseZap()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	defer logger.Sync()
	sugarLogger := logger.Sugar()
	sugarLogger.Debug("This is a debug message")
	sugarLogger.Info("This is an info message")
	sugarLogger.Warn("This is a warning message")
	sugarLogger.Error("This is an error message")
}
