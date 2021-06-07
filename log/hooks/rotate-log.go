package hooks

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type RotateFileConfig struct {
	filename   string
	maxSize    int
	maxBackups int
	maxAge     int
	level      logrus.Level
	formatter  logrus.Formatter
}

type RotateFileHook struct {
	config    RotateFileConfig
	logWriter io.Writer
}

func NewRotateFileConfig(filename string, maxsize, maxBackups, maxAge int) RotateFileConfig {
	return RotateFileConfig{
		filename:   filename,
		maxSize:    maxsize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		level:      logrus.DebugLevel,
		formatter:  &logrus.JSONFormatter{},
	}
}
func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	hook := RotateFileHook{
		config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.filename,
		MaxSize:    config.maxSize,
		MaxAge:     config.maxAge,
		MaxBackups: config.maxBackups,
	}
	return &hook, nil
}

func (r *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.config.level+1]
}

func (r *RotateFileHook) Fire(entry *logrus.Entry) error {
	b, err := r.config.formatter.Format(entry)
	if err != nil {
		return err
	}
	r.logWriter.Write(b)
	return nil
}
