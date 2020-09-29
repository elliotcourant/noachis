package logging

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.StandardLogger().SetLevel(logrus.TraceLevel)
	logrus.StandardLogger().Formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              true,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}
}

func NewLogger() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}
