package logger_test

import (
	"testing"

	"github.com/lwinmgmg/logger"
)

func TestNewLogging(t *testing.T) {
	logger := logger.NewLogging(logger.INFO, func() string { return "" }, 1000, logger.DEFAULT_PATTERN, logger.DEFAULT_TFORMAT, &logger.ConsoleWriter{})
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warning("Warning")
	logger.Error("Error")
	logger.Critical("Critical")
	logger.Close()
}

func TestDefaultLogging(t *testing.T) {
	logger := logger.DefaultLogging(logger.DEBUG)
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warning("Warning")
	logger.Error("Error")
	logger.Critical("Critical")
	logger.Close()
}
