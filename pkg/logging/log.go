package logging

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}
