package log_test

import (
	"os"
	"testing"

	"github.com/mrbitcoiner/log"
)

func TestExampleDefaultLog(t *testing.T) {
	log, _ := log.NewLog()
	log.Info("example: ", "default logger")
}

func TestExampleLogWithCustomWriter(t *testing.T) {
	log, _ := log.NewLog(log.WithWriter(os.Stdout))
	log.Info("logging to stdout")
}

func TestExampleLogWithCustomLevel(t *testing.T) {
	log, _ := log.NewLog(log.WithLevel(log.LOGERR))
	log.Err("error log")
}

func TestExampleLogWithCustomFileName(t *testing.T) {
	log, _ := log.NewLog(log.WithFilePath)
	log.Err("logging full file path")
}
