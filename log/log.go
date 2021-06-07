package log

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func NewEntry() *logrus.Entry {
	return logrus.NewEntry(Logger)
}
