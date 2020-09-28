package service

import (
	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/sirupsen/logrus"
)

type Service struct {
	closed int32
	log    *logrus.Entry
	db     engine.Engine
}
