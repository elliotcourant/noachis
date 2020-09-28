package exec

import (
	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/sirupsen/logrus"
)

type (
	Session struct {
		state     int32
		sessionId string
		log       *logrus.Entry
		dbEngine  engine.Engine

		dbTransaction engine.Transaction
	}
)
