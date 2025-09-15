package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Entry
}

func NewLogger(level string, appName string) *Logger {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = new(logrus.JSONFormatter)
	log.SetLevel(lvl)
	logger := Logger{
		logger: log.WithField("app", appName),
	}

	return &logger
}

func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		logger: l.logger.WithField(key, value),
	}
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	return &Logger{
		logger: l.logger.WithFields(fields),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.WithError(err),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger.WithContext(ctx),
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(args...)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *Logger) Warn(args ...any) {
	l.logger.Warn(args...)
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(args...)
}
