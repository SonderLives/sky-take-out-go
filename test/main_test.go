package test

import (
	"os"
	"testing"

	"goflow/internal/config"
	"goflow/internal/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.Init(&config.LogConfig{
		Mode:  "dev",
		Level: "debug",
	})
	os.Exit(m.Run())
}
