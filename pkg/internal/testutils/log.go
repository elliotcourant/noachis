package testutils

import (
	"testing"

	"github.com/elliotcourant/noachis/pkg/logging"
	"github.com/sirupsen/logrus"
)

func NewTestLogger(t *testing.T) *logrus.Entry {
	return logging.NewLogger().WithField("test", t.Name())
}
