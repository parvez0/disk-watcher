package watchers

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	baseLogger := logrus.New()
	logger := Logger{baseLogger}
	var err error
	level, exist := os.LookupEnv("LOGLEVEL")
	if !exist{
		level = "info"
	}
	logger.Level, err = logrus.ParseLevel(level)
	if err != nil{
		panic(err)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	return &logger
}
