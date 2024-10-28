package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(logFilePath, logFileName string) *logrus.Logger {
	logger := logrus.New()

	// Set output to standard output
	// logger.SetOutput(os.Stdout)

	// Set the formatter to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(logFilePath, 0777)
		if err != nil {
			panic(err)
		}
	}

	logFile, err := os.OpenFile(logFilePath+"/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(logFile)

	Log = logger
	return logger
}
