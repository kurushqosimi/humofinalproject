package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitializeLogger() (*logrus.Logger, error) {
	logger := logrus.New()
	logPath := "./pkg/logging/app.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	// Настройки логгера
	logger.SetOutput(file)
	return logger, nil
}
